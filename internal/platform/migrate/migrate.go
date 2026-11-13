package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Up runs SQL migrations from migrationsFS. dir is the root path inside the FS (often ".").
func Up(ctx context.Context, pool *pgxpool.Pool, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.UpContext(ctx, db, dir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

// Version returns current migration version (for tests).
func Version(db *sql.DB) (int64, error) {
	return goose.GetDBVersion(db)
}
