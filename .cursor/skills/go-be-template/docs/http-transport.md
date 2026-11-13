# HTTP

```
Client → chi → OpenAPI wrapper → transport.API → feature handler → httpx
```

## Middleware (`internal/http/router.go`)

Order: RequestID, Recoverer, request log, Timeout (55s), body limit (1 MiB).

Do not enable Chi `RealIP` until a trusted-proxy allowlist exists. It copies client `X-Forwarded-For` into `RemoteAddr`.

Parameter binding errors use `ErrorHandlerFunc` and `httpx.BadRequest`, so invalid path and query values are Problem JSON.

## Server timeouts

| Field | Value |
|-------|-------|
| ReadHeaderTimeout | 5s |
| ReadTimeout | 15s |
| WriteTimeout | 60s |
| IdleTimeout | 60s |
| Chi request timeout | 55s |

Write timeout stays above the Chi deadline.

## Handlers

Decode with `DisallowUnknownFields`. Validate with the generated `validate` tags. Call the service with `r.Context()`. Map sentinels to `httpx`. Log unexpected errors with `request_id` from chi middleware, then `httpx.InternalError`.

```go
h.log.Error("get user",
    slog.String("error", err.Error()),
    slog.String("request_id", chimw.GetReqID(r.Context())),
)
httpx.InternalError(w)
```

## Probes

| Route | Meaning |
|-------|---------|
| `GET /healthz` | Process is listening. No database call. |
| `GET /readyz` | `health.Checker.Ping`. 503 when the database is down or the process is shutting down. |

Feature handlers do not hold the pool for readiness.
