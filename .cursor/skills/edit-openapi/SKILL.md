---
name: edit-openapi
description: Edits HTTP routes contract-first with OpenAPI and oapi-codegen. Use when editing api/openapi.yaml, adding a route, changing a Problem response, or regenerating transport code.
paths:
  - "api/**"
  - "internal/transport/**/*.go"
---

# Edit the OpenAPI contract

`api/openapi.yaml` is the route inventory. Generated code is `internal/transport/openapi/api.gen.go`.

1. Edit the spec first: path, `operationId`, parameters, schemas, and `Problem` responses.
2. Run `make generate`.
3. Implement `transport.API` and the feature handler.

A handler bug needs no spec edit when the contract is unchanged. A new query parameter goes in the spec before the handler. Pagination reuses `Page` (1–1000000) and `PageSize` (1–100). UUID path params use `format: uuid` and the codegen name (`userId`). Request bodies set `additionalProperties: false` and `validate` tags in `x-oapi-codegen-extra-tags`. A breaking change is a new `/v2` path, after the user agrees.

Register routes with `openapi.HandlerWithOptions` and an `ErrorHandlerFunc` that writes Problem JSON through `httpx`. After a spec edit, `git diff --exit-code internal/transport/openapi` should be empty.
