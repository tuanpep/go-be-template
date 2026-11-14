---
name: place-packages
description: Places code in feature folders and keeps import direction. Use when creating a package, moving files, refactoring layers, or fixing an import cycle.
paths:
  - "internal/**/*.go"
  - "cmd/**/*.go"
---

# Place packages

| Path | Role |
|------|------|
| `cmd/api`, `cmd/migrate` | Process entry. No business logic. |
| `internal/app` | Composition root |
| `internal/<feature>/` | types, errors, repository, service, mapper, handler |
| `internal/platform` | postgres, httpx, health, log, migrate |
| `internal/transport` | OpenAPI server. Delegates to feature handlers. |
| `internal/store/db` | sqlc output. Leave it untouched. |

Import direction: `cmd` → `app` → `http`, `transport`, features → repository → `store/db`.

`platform` does not import a feature. Services do not import `store/db` or OpenAPI types. Handlers do not import `store/db`. A handler that needs data gets a repository method. A shared type goes in `platform` only when it is generic. Move an interface to its consumer to break a cycle.
