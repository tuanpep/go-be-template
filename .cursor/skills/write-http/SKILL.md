---
name: write-http
description: Writes Chi middleware, routing, and httpx responses. Use when editing the router, platform/httpx, health routes, timeouts, or an HTTP handler.
paths:
  - "internal/http/**/*.go"
  - "internal/platform/httpx/**/*.go"
  - "internal/transport/api.go"
  - "internal/**/handler.go"
---

# Write HTTP handlers

Client → chi → OpenAPI wrapper → `transport.API` → feature handler → `httpx`.

Middleware order in `internal/http/router.go`: RequestID, Recoverer, request log, Timeout (55s), body limit (1 MiB). Leave Chi `RealIP` off until a trusted-proxy allowlist exists.

Parameter errors use `openapi.HandlerWithOptions` and an `ErrorHandlerFunc` that writes Problem JSON through `httpx.BadRequest`.

| Timeout | Value |
|---------|-------|
| ReadHeaderTimeout | 5s |
| ReadTimeout | 15s |
| WriteTimeout | 60s |
| IdleTimeout | 60s |
| Chi request timeout | 55s |

Write timeout stays above the Chi deadline.

Handlers decode with `DisallowUnknownFields`, validate the generated `validate` tags, call the service with `r.Context()`, and map sentinels to `httpx`. `httpx.ParsePage` rejects an out-of-range page. Log unexpected errors with `request_id` from `chimw.GetReqID`, then `httpx.InternalError`.

`GET /healthz` means the process is listening and does not touch the database. `GET /readyz` calls `health.Checker.Ping` and returns 503 when the database is down or the process is shutting down.
