package slog

import "testing"

func TestTHandler(t *testing.T) {
	log := TestLogger(t)
	log.Debug("hello", "test", t.Name())
}
