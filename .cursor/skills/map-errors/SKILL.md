---
name: map-errors
description: >-
  Maps feature sentinel errors to slog and RFC 7807 Problem JSON via httpx.
  Use when editing internal/**/errors.go, handlers, or internal/platform/httpx.
paths:
  - "internal/**/errors.go"
  - "internal/**/handler.go"
  - "internal/platform/httpx/**/*.go"
---

# Map errors

Sentinels in feature `errors.go`; use `errors.Is`.

## Instructions

| Layer | Rule |
|-------|------|
| Repository | Driver → sentinel; else `%w` |
| Service | Return sentinels |
| Handler | One helper mapping to `httpx.NotFound` / `Conflict` / `BadRequest` / log + `InternalError` |

No credentials in logs. `r.Context()` through service and repository.