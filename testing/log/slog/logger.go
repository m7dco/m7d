package slog

import (
	"context"
	"log/slog"
	"testing"
)

type THandler struct {
	test *testing.T
}

func (t *THandler) Enabled(ctx context.Context, l slog.Level) bool {
	return true
}

func (t *THandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]string, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a.String())
		return true
	})

	t.test.Log(r.Time, r.Level, r.Message, attrs)
	return nil
}

func (t *THandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return t
}

func (t *THandler) WithGroup(name string) slog.Handler {
	return t
}

func TestLogger(t *testing.T) *slog.Logger {
	return slog.New(&THandler{t})
}
