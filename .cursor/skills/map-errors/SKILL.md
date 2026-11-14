---
name: map-errors
description: Maps sentinel errors to slog and Problem responses. Use when editing errors.go, handler error branches, or platform/log.
paths:
  - "internal/**/errors.go"
  - "internal/**/handler.go"
  - "internal/platform/log/**/*.go"
  - "internal/platform/httpx/**/*.go"
---

# Map errors

Sentinels live in the feature `errors.go`. Compare them with `errors.Is`.

| Layer | Style |
|-------|-------|
| Repository | Map driver codes to sentinels. Wrap other I/O with `%w`. |
| Service | Return sentinels. Wrap with `%w` only when `errors.Is` still works. |
| Handler | `errors.Is` → `httpx` status. |

| Outcome | Response |
|---------|----------|
| Not found | `httpx.NotFound` |
| Conflict | `httpx.Conflict` |
| Bad input | `httpx.BadRequest` with `httpx.ValidationDetail` |
| Unexpected | Log, then `httpx.InternalError` with no driver text |

Error logs include `request_id` from `chimw.GetReqID`. Leave passwords, tokens, and credential URLs out of logs. Pass `r.Context()` into the service and the repository. Chi `Recoverer` handles handler panics.
