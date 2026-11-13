# API contract

`api/openapi.yaml` is the route inventory. Generated code is `internal/transport/openapi/api.gen.go`.

## Workflow

1. Edit the spec: path, `operationId`, parameters, schemas, and Problem responses.
2. `make generate` (`go tool oapi-codegen`).
3. Implement `transport.API` and the feature handler.

A handler-only bug does not need a spec change when the contract is unchanged.

## Decisions

| Situation | Do |
|-----------|----|
| New query param | Add it under `parameters` first. |
| Error body | `components/schemas/Problem` via `httpx`. |
| Pagination | Reuse `Page` (1–1000000) and `PageSize` (1–100). |
| UUID path | `format: uuid`. Name matches codegen (`userId`). |
| Request body | `additionalProperties: false`. Put `validate` tags in `x-oapi-codegen-extra-tags` so max length is enforced. |
| Breaking change | New `/v2` path, with the user's approval. |

CI fails when `internal/transport/openapi` drifts from the spec.

Register routes with `openapi.HandlerWithOptions` and an `ErrorHandlerFunc` that writes Problem JSON. Do not use the generated `http.Error` default.
