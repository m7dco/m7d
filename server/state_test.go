package server

import (
	"strconv"
	"testing"
	"time"
)

func TestStateReady(t *testing.T) {
	s := newState()

	if s.IsRunning() || s.IsReady() {
		t.Fatal("wrong state; expected !running && !ready")
	}

	rf := s.ReadyReporter(t.Name())
	rf(true)

	if !s.IsRunning() || !s.IsReady() {
		t.Fatal("wrong state; expected running && ready")
	}

	rf2 := s.ReadyReporter(t.Name() + "2")
	rf2(false)

	if !s.IsRunning() || s.IsReady() {
		t.Fatal("wrong state; expected running && !ready")
	}
}

func TestStateHealthy(t *testing.T) {
	s := newState()

	if !s.IsHealthy() {
		t.Fatal("wrong state; expected healthy")
	}
}

func TestRunning(t *testing.T) {
	s := newState()
	running := false

	go func() {
		rf := s.ReadyReporter(t.Name())
		rf(true)
	}()

	select {
	case <-s.Running():
		running = true

	case <-time.After(900 * time.Millisecond):
	}

	if !running {
		t.Fatal("running never signalled")
	}
}

func TestStateDetectRaceCond(t *testing.T) {
	s := newState()
	for i := 0; i < 20; i++ {
		go func(i int) {
			rr := s.ReadyReporter(t.Name() + "ready" + strconv.Itoa(i))
			hr := s.HealthyReporter(t.Name() + "healthy" + strconv.Itoa(i))

			rr(true)
			hr(false)
			rr(false)
			hr(true)
		}(i)
	}

	go func() {
		<-s.Running()
	}()

	time.Sleep(100 * time.Millisecond)
}
