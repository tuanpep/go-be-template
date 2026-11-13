# Structure

## Layout

| Path | Role |
|------|------|
| `cmd/api`, `cmd/migrate` | Process entry. No business logic. |
| `internal/app` | Composition root |
| `internal/<feature>/` | types, errors, repository, service, mapper, handler |
| `internal/platform` | postgres, httpx, health, log, migrate. No feature imports. |
| `internal/transport` | OpenAPI server. Delegates to feature handlers. |
| `internal/store/db` | sqlc output. Never hand-edit. |

Golden path: `internal/user`.

## Decisions

| Situation | Do |
|-----------|----|
| New capability | New `internal/<feature>/`. Not `internal/services/`. |
| Second feature | Its own handler. Extend `transport.NewAPI`. |
| Shared type | `internal/platform` only when it is generic. Otherwise a small DTO per feature. |
| Handler needs data | Add a repository method. Do not import `store/db`. |
| Service needs a response type | Map in `mapper.go`. Service returns the domain type. |
| Import cycle | Move the interface to the consumer. Do not add `internal/common`. |
| Global `entity/` tree | Keep feature folders unless the user explicitly insists. |

## Imports

```
cmd → app → http, transport, features → repository → store/db
```

`platform` does not import `internal/<feature>`.
Services do not import `store/db` or OpenAPI types.
Handlers do not import `store/db`.
