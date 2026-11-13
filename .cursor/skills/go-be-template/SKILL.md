---
name: go-be-template
description: Routes work in this Chi, sqlc, pgx, OpenAPI, and goose template. Use when editing the repository, scaffolding a feature, reviewing a diff, or when the user mentions go-be-template.
paths:
  - "**/*.go"
  - "**/*.sql"
  - "**/openapi.yaml"
  - "**/sqlc.yaml"
  - "Makefile"
  - "Dockerfile"
  - "docker-compose.yml"
  - ".env.example"
  - ".golangci.yml"
  - ".github/**"
---

# go-be-template

Golden path: `internal/user`. Copy that shape for a new feature.

## Do this

1. Pick one row in the task router and read that doc before writing code.
2. Follow the edit order and the principles below.
3. Run the row's verify script. If the change crosses layers, run `scripts/validate.sh`.

If the user asks for something that breaks a principle, explain the conflict and wait for an explicit yes.

## Task router

Paths are relative to `.cursor/skills/go-be-template/`.

| Task | Doc | Verify |
|------|-----|--------|
| Packages, imports | [docs/structure.md](docs/structure.md) | `scripts/verify-structure.sh` |
| OpenAPI, codegen | [docs/api-contract.md](docs/api-contract.md) | `scripts/verify-openapi.sh` |
| Chi, handlers, httpx | [docs/http-transport.md](docs/http-transport.md) | `scripts/verify-http.sh` |
| `internal/app`, cmd | [docs/application-wiring.md](docs/application-wiring.md) | `scripts/verify-wiring.sh` |
| SQL, migrations, repositories | [docs/database.md](docs/database.md) | `scripts/verify-data.sh` |
| Tests, CI | [docs/testing-ci.md](docs/testing-ci.md) | `scripts/verify-testing.sh` |
| Docker, env, probes | [docs/operations.md](docs/operations.md) | `scripts/verify-operations.sh` |
| Errors, logs | [docs/errors-logging.md](docs/errors-logging.md) | `scripts/verify-errors.sh` |
| Secrets, SQL safety | [docs/security.md](docs/security.md) | `scripts/verify-security.sh` |
| Review or audit | [docs/code-review.md](docs/code-review.md) | `scripts/verify-review.sh` |
| New feature end to end | [examples.md](examples.md) | `scripts/validate.sh` |

## Edit order

| Change | Order |
|--------|-------|
| New HTTP API | `api/openapi.yaml` → `make generate` → handler |
| New persistence | migration → `sql/schema` → `sql/queries` → `make generate` → repository |
| New feature | OpenAPI and data → `internal/<feature>/` → `app.go` → `transport/api.go` |

Never hand-edit `internal/store/db/*` or `internal/transport/openapi/api.gen.go`.

## Principles

1. Thin `cmd/*/main.go`. Wire dependencies only in `internal/app`.
2. Feature packages are `internal/<name>/`, not global layer folders.
3. Call chain is `transport.API` → handler → service → repository → sqlc.
4. Domain types live in the feature. Services do not return `db.*` or OpenAPI types.
5. Change `api/openapi.yaml` before handlers, then `make generate`.
6. Migrations run from `cmd/migrate`, not from API startup.
7. Chi and `platform/httpx`. Readiness goes through `health.Checker`.
8. Config is environment only. `DATABASE_URL` has no default. No secrets in the repo.
9. `/healthz` does not ping the database. `/readyz` does.
10. Smallest correct diff.

## Done

```bash
.cursor/skills/go-be-template/scripts/validate.sh
```

## Example

User: "Add a widgets API." Read [examples.md](examples.md), implement, run `scripts/validate.sh`.
