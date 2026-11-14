---
name: add-persistence
description: >-
  Adds goose migrations, sqlc queries, and pgx repositories per this template.
  Use when editing migrations/, sql/, sqlc.yaml, internal/store, or feature
  repository.go.
paths:
  - "migrations/**"
  - "sql/**"
  - "sqlc.yaml"
  - "internal/**/repository.go"
  - "internal/store/**"
---

# Add persistence

Pool: `internal/platform/postgres`. Migrations: `migrations/`. sqlc: `sql/schema` + `sql/queries` → `internal/store/db`.

## Instructions

1. `migrations/NNNNN_*.sql` (Up/Down) — one change per file.
2. Mirror schema/indexes in `sql/schema/schema.sql`.
3. `sql/queries/<feature>.sql` → `make generate`.
4. `Repository` interface in `service.go`; `SQLRepository` maps rows in `repository.go`.

`pgx.ErrNoRows` → `ErrNotFound`. Unique `23505` → conflict sentinel. List + count in one read-only tx. Multi-step write → tx in repository. Services: no `pgx`/`store/db`.