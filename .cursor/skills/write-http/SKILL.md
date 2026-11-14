---
name: write-http
description: >-
  Implements Chi routing, httpx Problem responses, and feature HTTP handlers in
  this template. Use when editing internal/http, internal/transport/api.go,
  platform/httpx, or feature handler.go.
paths:
  - "internal/http/**/*.go"
  - "internal/platform/httpx/**/*.go"
  - "internal/transport/api.go"
  - "internal/**/handler.go"
---

# Write HTTP

Chain: Chi → OpenAPI wrapper → `transport.API` → feature handler → `httpx`.

## Instructions

`internal/http/router.go`: RequestID, Recoverer, log, Timeout 55s, body 1 MiB. No `RealIP` without trusted proxy allowlist.

Server timeouts (`internal/http/timeouts.go`): ReadHeader 5s, Read 15s, Write 60s, Idle 60s; Write > Chi timeout.

Handlers: `Service` interface; `DisallowUnknownFields`; validate tags; `r.Context()` into service; one sentinel mapper (**map-errors**); log with `request_id`; internal → `httpx.InternalError`.

`/healthz` and `/readyz` are OpenAPI routes (not feature handlers). Probe meaning and deploy order: **configure-runtime**.