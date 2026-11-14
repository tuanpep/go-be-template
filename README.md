# go-be-template

Production-oriented Go HTTP API template (Chi, sqlc, pgx, OpenAPI).

## Stack

| Layer | Choice |
|-------|--------|
| HTTP | Chi on `net/http` |
| Contract | OpenAPI 3 + oapi-codegen |
| Database | PostgreSQL, pgx pool, sqlc, goose |
| Structure | `internal/app` composition, feature package with a public `Service` interface |
| Logging | `log/slog` |

## Layout

```
cmd/api/          # thin main → internal/app
cmd/migrate/      # database migrations only (run before deploy / locally)
internal/app/     # composition root (wire dependencies)
internal/user/    # sample feature: Service interface, handler and SQL at the edges
internal/platform/# postgres, httpx, health, migrate helpers
internal/transport/ # OpenAPI facade (delegates to features)
integration/      # HTTP + DB tests (needs DATABASE_URL)
```

## Quick start

`/v1/users` is a public sample. This template does not authenticate callers.

`DATABASE_URL` is required. The API exits when it is unset. `make run` and `make migrate` load `.env` if you created one; otherwise they use the compose URL.

```bash
cp .env.example .env
make run
```

```bash
curl -s localhost:8080/healthz
curl -s -X POST localhost:8080/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","name":"Alice"}'
```

## Docker

```bash
make docker-build
# Run migrations as a one-off job, then start API:
docker run --rm -e DATABASE_URL=... go-be-template:local /app/migrate
docker run -p 8080:8080 -e DATABASE_URL=... go-be-template:local
```

Migrations are **not** applied automatically when the API starts (matches k8s Job / init-container patterns).

## Development

- `make generate` — sqlc + oapi-codegen (CI enforces no drift)
- `make lint` — golangci-lint
- `INTEGRATION`: CI sets `DATABASE_URL` and runs `go test ./integration/...`
- HTTP: shared timeouts in `internal/http/timeouts.go`; request bodies capped at 1 MiB

## Add a feature

1. Extend `api/openapi.yaml` → `make generate`
2. Add `internal/<feature>/` copying `internal/user`: `Service` interface, use-case inputs, handler and SQL repository at the edges
3. Wire in `internal/app/app.go` and delegate from `internal/transport/api.go`

See `.cursor/skills/README.md` for agent skills (architecture, scaffold, grow, and task skills).

## AI agents

- Start: `AGENTS.md` and [.cursor/skills/README.md](.cursor/skills/README.md)
- New feature: **scaffold-feature**; growth: **grow-feature**; layout: **codebase-architecture**
- After a spec or SQL change: `make generate`, then `git diff --exit-code internal/store/db internal/transport/openapi`
