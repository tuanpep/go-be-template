---
name: scaffold-feature
description: >-
  Scaffolds a new Chi, sqlc, pgx, OpenAPI, and goose resource in this template
  (copy internal/user). Use when adding a feature, resource, or endpoint set.
paths:
  - "api/**"
  - "sql/**"
  - "migrations/**"
  - "internal/**"
  - "integration/**"
---

# Scaffold a feature

Copy `internal/user` (list/create/get/patch/delete; list+count in one read-only tx). Scale later: **grow-feature**. Layout rules: **codebase-architecture**.

## Instructions

Run in order; do not re-derive steps from other skills—follow each linked skill for that layer only.

1. **edit-openapi** — contract and `make generate`.
2. **add-persistence** — migrations, schema mirror, sqlc queries, repository.
3. **Package** — `internal/<feature>/`: `types.go`, `errors.go`, `service.go` (`Service` + `Repository`), `repository.go`, `mapper.go`, `handler.go`. `NewService` returns `Service`; handler takes `Service`; `CreateInput`/`UpdateInput` structs.
4. **wire-app** — construct in `internal/app`; delegate in `internal/transport/api.go`.
5. **map-errors** — sentinels + handler mapper.
6. **write-http** — handler HTTP behavior (if not already aligned with `internal/user`).
7. **write-tests** — `integration/<feature>_api_test.go` for new HTTP surfaces (skip without `DATABASE_URL`).

Migrations via `cmd/migrate` only (**configure-runtime**).