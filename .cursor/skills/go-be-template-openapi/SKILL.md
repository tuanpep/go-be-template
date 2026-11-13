---
name: go-be-template-openapi
description: Applies contract-first OpenAPI and oapi-codegen for go-be-template. Use when editing api/openapi.yaml, HTTP routes, Problem responses, or generated transport code.
paths:
  - "api/**"
  - "internal/transport/**/*.go"
---

# OpenAPI

1. Read [api-contract.md](../go-be-template/docs/api-contract.md).
2. Edit `api/openapi.yaml`, then `make generate`, then handlers.
3. Run `.cursor/skills/go-be-template/scripts/verify-openapi.sh`.
