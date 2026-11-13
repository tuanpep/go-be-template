# Database

| Tool | Role |
|------|------|
| pgx pool | `internal/platform/postgres`. MaxConns 20, MinConns 2. No startup ping. |
| goose | `migrations/` via `cmd/migrate` and `migrations/embed.go` |
| sqlc | `sql/schema` + `sql/queries` → `internal/store/db` |

## Workflow

1. `migrations/NNNNN_description.sql` with `-- +goose Up` and `-- +goose Down`.
2. Mirror tables and list indexes in `sql/schema/schema.sql`.
3. Add `sql/queries/<feature>.sql`.
4. `make generate`.
5. Map `db` rows to domain types in the repository.

One logical change per migration. Apply with `make migrate`, not from the API.

## Decisions

| Situation | Do |
|-----------|----|
| Not found | `pgx.ErrNoRows` → `ErrNotFound` |
| Duplicate email | SQLSTATE `23505` and constraint `users_email_key` → `ErrDuplicateEmail` |
| Other unique violation | Wrap and return it. Do not call every `23505` a duplicate email. |
| List | `ORDER BY created_at DESC, id DESC`, `LIMIT`/`OFFSET`, plus a count query. |
| Page bounds | Reject out-of-range pages in `httpx.ParsePage`. Do not clamp into a negative `int32`. |
| Indexes | A `UNIQUE` column already has an index. Add an index for the list order, not a second one on `email`. |
| Multi-step write | Transaction inside the repository. |

Services do not import `pgx` or `internal/store/db`.
