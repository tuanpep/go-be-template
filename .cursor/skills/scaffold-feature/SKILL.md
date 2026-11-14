---
name: scaffold-feature
description: Scaffolds a resource in this Chi, sqlc, pgx, OpenAPI, and goose template. Use when adding a feature or a new endpoint set.
---

# Scaffold a feature

Copy `internal/user`. That package is the golden path: list, create, get, partial update, and delete. List and count share one read-only transaction.

## Order

1. Add the paths to `api/openapi.yaml` with `operationId`, `Problem` errors, `additionalProperties: false`, and `x-oapi-codegen-extra-tags` `validate` tags. Run `make generate`.
2. Add `migrations/NNNNN_name.sql` using the next free number, with goose Up and Down. Mirror the table and list indexes in `sql/schema/schema.sql`. Add `sql/queries/<feature>.sql`. List queries order by a timestamp and a unique id. Run `make generate` again.
3. Add `internal/<feature>/` with `types.go`, `errors.go`, `repository.go`, `service.go`, `mapper.go`, and `handler.go`.
4. Construct repository, service, and handler in `internal/app`. Delegate from `internal/transport/api.go`.

## Files

| File | Role |
|------|------|
| `types.go` | Domain struct |
| `errors.go` | Sentinels such as `ErrNotFound` |
| `repository.go` | `Repository` and `SQLRepository`. Map `db` rows to domain types here. |
| `service.go` | Use cases. Domain types only. |
| `mapper.go` | Domain to OpenAPI |
| `handler.go` | Decode, validate, call the service, respond with `httpx` |

## Rules

- `cmd/*/main.go` stays thin. Wire dependencies only in `internal/app`.
- New behavior goes in `internal/<feature>/`.
- Call chain is `transport.API` → handler → service → repository → sqlc.
- Leave `internal/store/db/` and `internal/transport/openapi/api.gen.go` untouched. Regenerate with `make generate`.
- Migrations run from `cmd/migrate`.
- `DATABASE_URL` is required and has no default. `/healthz` does not ping the database. `/readyz` does.
- Integration tests skip when `DATABASE_URL` is unset and use `uuid.NewString()` in unique fields.

## Check

`make generate`, then `git diff --exit-code internal/store/db internal/transport/openapi`, then `go test` for the packages you changed.
