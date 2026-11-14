---
name: edit-openapi
description: >-
  Edits contract-first OpenAPI and oapi-codegen output for this template. Use
  when changing api/openapi.yaml, generated transport types, routes, or
  operationId wiring.
paths:
  - "api/**"
  - "internal/transport/**/*.go"
---

# Edit OpenAPI

## Instructions

1. Edit `api/openapi.yaml` (path, `operationId`, schemas, Problem responses).
2. `make generate`.
3. Delegate in `internal/transport/api.go`; implement in feature `handler`.

Bodies: `additionalProperties: false`, `x-oapi-codegen-extra-tags` `validate`. Pagination: `Page` 1–1000000, `PageSize` 1–100. UUID params: `format: uuid`. Breaking changes → new `/v2` after user agrees.

Router uses `openapi.HandlerWithOptions` + `ErrorHandlerFunc` → `httpx` Problem JSON.