package server

import (
	"context"
	"flag"
	"runtime"
	"testing"
	"time"

	"github.com/m7dco/m7d"
	"github.com/m7dco/m7d/env"
)

func checkRRH(t *testing.T, h *Host) {
	t.Logf("server:Host; running:%+v ready:%+v healthy:%+v", h.State.IsRunning(), h.State.IsReady(), h.State.IsHealthy())

	if !h.State.IsRunning() || !h.State.IsReady() || !h.State.IsHealthy() {
		t.Fatal()
	}
}

var (
	canonicalMainTimeout = flag.Duration("canonical_main_timeout", 100*time.Millisecond, "Timeout for the main")
)

func TestCanonicalMain(t *testing.T) {
	e := m7d.Check(env.Init())
	host := NewHost(e, 8080, 8081)
	t.Log(host)

	if host.State.IsRunning() {
		t.Fatal("should not be running yet")
	}

	ctx, cancel := context.WithCancel(context.Background())
	res := make(chan error)
	go func() {
		res <- host.Run(ctx)
	}()

	runtime.Gosched()
	<-host.State.Running()

	checkRRH(t, host)

	time.Sleep(*canonicalMainTimeout)

	cancel()

	runErr := <-res
	t.Log(runErr)
}
