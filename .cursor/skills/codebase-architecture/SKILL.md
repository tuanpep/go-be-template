---
name: codebase-architecture
description: >-
  Describes repo layout, import boundaries, and where code belongs in this Go
  backend template. OSS comparisons live in references/oss-repos.md. Use when
  onboarding, moving packages, fixing import cycles, or deciding where new code
  goes.
paths:
  - "go.mod"
  - "AGENTS.md"
  - "README.md"
  - "cmd/**"
  - "internal/app/**"
---

# Codebase architecture

Modular monolith: one `go.mod`, one database, `cmd/api` + `cmd/migrate`. OSS trees: [references/oss-repos.md](references/oss-repos.md).

New resource → **scaffold-feature**. Scaling → **grow-feature**. Layer edits → task skills in [.cursor/skills/README.md](../README.md).

## Instructions

### Request flow

`openapi.yaml` → `internal/http` → `internal/transport` (delegate) → feature `handler` → `Service` → `repository` → `internal/store/db` → PostgreSQL.

Only `internal/app` constructs dependencies. `cmd/*` is thin config + `app.Run`.

### Where code lives

| Path | Role |
|------|------|
| `api/` | Contract; `make generate` |
| `sql/`, `migrations/` | sqlc + goose |
| `cmd/api`, `cmd/migrate` | Entry only |
| `internal/app` | Wiring |
| `internal/<feature>/` | `Service`, handler, repository (copy `internal/user`) |
| `internal/platform/` | postgres, httpx, health, log, migrate helpers |
| `internal/http/` | Chi router, timeouts, middleware |
| `internal/transport/` | Generated `ServerInterface`; forward to handlers |
| `internal/store/db/` | sqlc output; do not edit |
| `integration/` | HTTP + DB tests |

Feature exports: `Service`, domain types, sentinels. Keep OpenAPI and sqlc types on the handler/repository edges.

### Imports

`cmd` → `internal/app` → `internal/http` | `internal/transport` | features → repository → `internal/store/db`.

- `platform` does not import features.
- Services: no `store/db`, no OpenAPI types.
- Handlers: no `store/db`.
- Features call other features via `Service`, not another repository.
- Shared types belong in `platform` only when generic.
- Put the interface on its consumer to break cycles.

### Do not

- Global layer packages (`internal/handler`, `internal/repository` for all features).
- Logic in `transport` or wiring outside `app`.
- Routes absent from `api/openapi.yaml`.
- Migrate on API startup.

### PR checklist

Contract/SQL updated and generated trees clean (`make generate`; see AGENTS.md). Wiring only in `app`. Growth follows **grow-feature** stage order.