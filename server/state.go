package server

import "sync"

type ReadyFunc func(isReady bool)
type HealthyFunc func(isHealthy bool)

type rstate struct {
	ready          bool
	reportsReady   bool
	healthy        bool
	reportsHealthy bool
}

type State struct {
	rwm       *sync.RWMutex
	running   bool
	rchan     chan any
	reporters map[string]*rstate
}

func newState() *State {
	return &State{&sync.RWMutex{}, false, make(chan any), map[string]*rstate{}}
}

func (s *State) read(rf func()) {
	s.rwm.RLock()
	defer s.rwm.RUnlock()

	rf()
}

func (s *State) mutate(mf func()) {
	s.rwm.Lock()
	defer s.rwm.Unlock()

	mf()
}

func (s *State) reporter(name string) *rstate {
	r, ok := s.reporters[name]
	if !ok {
		r = &rstate{}
		s.reporters[name] = r
	}
	return r
}

func (s *State) Running() <-chan any {
	return s.rchan
}

func (s *State) ReadyReporter(name string) ReadyFunc {
	var rf ReadyFunc
	s.mutate(func() {
		r := s.reporter(name)
		r.reportsReady = true

		rf = func(isReady bool) {
			s.mutate(func() {
				r.ready = isReady
				if !isReady || s.running {
					return
				}

				s.running = true
				close(s.rchan)
			})
		}
	})
	return rf
}

func (s *State) HealthyReporter(name string) HealthyFunc {
	var hf HealthyFunc
	s.mutate(func() {
		r := s.reporter(name)
		r.reportsHealthy = true

		hf = func(isHealthy bool) {
			s.mutate(func() {
				r.healthy = isHealthy
			})
		}
	})
	return hf
}

func (s *State) IsReady() bool {
	ready := false
	s.read(func() {
		if !s.running {
			return
		}

		for _, r := range s.reporters {
			if r.reportsReady && !r.ready {
				return
			}
		}

		ready = true
	})

	return ready
}

func (s *State) IsHealthy() bool {
	healthy := true
	s.read(func() {
		for _, r := range s.reporters {
			if r.reportsHealthy && !r.healthy {
				healthy = false
				return
			}
		}
	})
	return healthy
}

func (s *State) IsRunning() bool {
	r := false
	s.read(func() {
		r = s.running
	})
	return r
}
