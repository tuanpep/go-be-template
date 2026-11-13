package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tuanbt/go-be-template/internal/platform/health"
)

var _ health.Checker = (*PoolChecker)(nil)

// PoolChecker adapts a pgx pool for readiness probes.
type PoolChecker struct {
	pool *pgxpool.Pool
}

func NewPoolChecker(pool *pgxpool.Pool) *PoolChecker {
	return &PoolChecker{pool: pool}
}

func (c *PoolChecker) Ping(ctx context.Context) error {
	return Ping(ctx, c.pool)
}
