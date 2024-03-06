package server

import (
	"context"
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

func TestCanonicalMain(t *testing.T) {
	e := m7d.Check(env.Init())
	host := NewHost(e, 8080, 8081)
	t.Log(host)

	if host.State.IsRunning() {
		t.Fatal("should not be running yet")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		host.Run(ctx)
	}()

	runtime.Gosched()
	<-host.State.Running()

	checkRRH(t, host)

	time.Sleep(1000 * time.Second)

	cancel()
}
