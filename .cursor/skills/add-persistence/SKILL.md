---
name: add-persistence
description: Adds goose migrations, sqlc queries, and pgx repositories. Use when writing migrations, sql/schema, sql/queries, sqlc.yaml, or repository code.
paths:
  - "migrations/**"
  - "sql/**"
  - "sqlc.yaml"
  - "internal/**/repository.go"
  - "internal/store/**"
---

# Add persistence

pgx pool lives in `internal/platform/postgres` (MaxConns 20, MinConns 2, no startup ping). goose uses `migrations/` and `migrations/embed.go`. sqlc reads `sql/schema` and `sql/queries` and writes `internal/store/db`.

1. Add `migrations/NNNNN_description.sql` with goose Up and Down. One logical change per file.
2. Mirror tables and list indexes in `sql/schema/schema.sql`.
3. Add `sql/queries/<feature>.sql`.
4. Run `make generate`.
5. Map `db` rows to domain types in the repository. Apply with `make migrate`.

`pgx.ErrNoRows` becomes `ErrNotFound`. SQLSTATE `23505` on constraint `users_email_key` becomes `ErrDuplicateEmail`. Any other unique violation is wrapped and returned. Lists use `ORDER BY created_at DESC, id DESC`, `LIMIT`/`OFFSET`, and a count query inside one read-only transaction. A `UNIQUE` column is already indexed. Add an index for the list order. A multi-step write is a transaction inside the repository. Services do not import `pgx` or `internal/store/db`.

After a SQL edit, `git diff --exit-code internal/store/db` should be empty.
