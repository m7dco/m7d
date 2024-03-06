package server

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

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

func (h *Host) runServer(cancel func(), name string, srv *http.Server) {
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return

	}

	h.log.Error("server stopped", "name", name, "err", err)
	cancel()
}

func (h *Host) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	go h.runServer(cancel, "external", h.Server)
	go h.runServer(cancel, "internal", h.InternalServer)

	h.ready(true)

	<-ctx.Done()
}
