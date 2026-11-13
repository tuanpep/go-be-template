# Errors and logging

| Layer | Style |
|-------|-------|
| Repository | Map driver codes to sentinels. Wrap other I/O with `%w`. |
| Service | Return sentinels. Wrap with `%w` only when `errors.Is` still works. |
| Handler | `errors.Is` → `httpx` status. |

Sentinels live in the feature `errors.go`. Compare with `errors.Is`, not string text.

| Outcome | Response |
|---------|----------|
| Not found | `httpx.NotFound` |
| Conflict | `httpx.Conflict` |
| Bad input | `httpx.BadRequest` with `httpx.ValidationDetail` |
| Unexpected | log, then `httpx.InternalError` with no driver text |

Error logs include `request_id` (`chimw.GetReqID`). Do not log passwords, tokens, or raw credential URLs.

Do not panic for control flow. Chi `Recoverer` handles handler panics.

Pass `r.Context()` into the service and the repository.
