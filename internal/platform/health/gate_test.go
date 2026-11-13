package health

import (
	"context"
	"errors"
	"testing"
)

type okChecker struct{}

func (okChecker) Ping(context.Context) error { return nil }

func TestGate_CloseFailsReadiness(t *testing.T) {
	t.Parallel()
	g := NewGate(okChecker{})
	if err := g.Ping(context.Background()); err != nil {
		t.Fatalf("open: %v", err)
	}
	g.Close()
	if err := g.Ping(context.Background()); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("closed: got %v", err)
	}
}
