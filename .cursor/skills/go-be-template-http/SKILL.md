---
name: go-be-template-http
description: Applies Chi middleware, routing, and httpx for go-be-template. Use when editing the router, platform/httpx, health routes, or an HTTP handler.
paths:
  - "internal/http/**/*.go"
  - "internal/platform/httpx/**/*.go"
  - "internal/transport/api.go"
  - "internal/**/handler.go"
---

# HTTP

1. Read [http-transport.md](../go-be-template/docs/http-transport.md).
2. Handlers decode, validate, call the service, and respond with `httpx`.
3. Run `.cursor/skills/go-be-template/scripts/verify-http.sh`.
