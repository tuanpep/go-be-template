package health

import "context"

// Checker reports dependency readiness (database, cache, etc.).
type Checker interface {
	Ping(ctx context.Context) error
}
