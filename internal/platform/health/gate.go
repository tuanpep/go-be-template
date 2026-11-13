package health

import (
	"context"
	"errors"
	"sync/atomic"
)

// ErrUnavailable means the process is not accepting traffic (shutting down).
var ErrUnavailable = errors.New("not accepting traffic")

// Gate wraps a dependency checker and can fail readiness before shutdown drains requests.
type Gate struct {
	next Checker
	open atomic.Bool
}

func NewGate(next Checker) *Gate {
	g := &Gate{next: next}
	g.open.Store(true)
	return g
}

func (g *Gate) Ping(ctx context.Context) error {
	if !g.open.Load() {
		return ErrUnavailable
	}
	return g.next.Ping(ctx)
}

// Close marks the process unready. In-flight requests keep using the pool.
func (g *Gate) Close() {
	g.open.Store(false)
}
