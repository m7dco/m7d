package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/m7dco/m7d/env"
	xhttp "github.com/m7dco/m7d/server/net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Host struct {
	env            *env.Env
	log            *slog.Logger
	State          *State
	ready          ReadyFunc
	Server         *http.Server
	Mux            *http.ServeMux
	InternalServer *http.Server
	InternalMux    *http.ServeMux
	PRegistry      *prometheus.Registry
}

func newHttpServer(port int) (*http.Server, *http.ServeMux) {
	server := &http.Server{}
	server.Addr = ":" + strconv.Itoa(port)
	mux := http.NewServeMux()
	server.Handler = mux
	return server, mux
}

func NewHost(e *env.Env, port, internalPort int) *Host {
	state := newState()
	server, mux := newHttpServer(port)
	internalServer, internalMux := newHttpServer(internalPort)
	internalMux.Handle("/started", &xhttp.ProbezHandler{state.IsRunning, "started"})
	internalMux.Handle("/ready", &xhttp.ProbezHandler{state.IsReady, "ready"})
	internalMux.Handle("/healthy", &xhttp.ProbezHandler{state.IsHealthy, "healthy"})

	reg := prometheus.NewRegistry()
	internalMux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	registerDefaultMetrics(reg)

	log := e.Log.WithGroup("server.Host")

	ready := state.ReadyReporter("host")

	h := &Host{e, log, state, ready, server, mux, internalServer, internalMux, reg}
	return h
}

func registerDefaultMetrics(reg *prometheus.Registry) {
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
	))
}

func (h *Host) runServer(errc chan error, name string, srv *http.Server) {
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return
	}

	h.log.Error("server stopped", "name", name, "err", err)
	errc <- err
}

func (h *Host) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	errc := make(chan error, 2)
	go h.runServer(errc, "external", h.Server)
	go h.runServer(errc, "internal", h.InternalServer)

	h.ready(true)

	res := []error{}

	select {
	case <-ctx.Done():
		h.Server.Close()
		h.InternalServer.Close()

	case err := <-errc:
		res = append(res, err)
		h.Server.Close()
		h.InternalServer.Close()
	}

	cancel()

	shutdownAt := time.After(100 * time.Millisecond)
	for {
		select {
		case <-shutdownAt:
			return errors.Join(res...)

		case err := <-errc:
			res = append(res, err)
		}
	}
}
