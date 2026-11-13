# Add a feature

Copy `internal/user`. Replace `widget` with the resource name. Use the next migration number, not a number that already exists.

## 1. Spec

Add `/v1/widgets` to `api/openapi.yaml` with `operationId`s `listWidgets`, `createWidget`, and `getWidget`. Schemas: `Widget`, `CreateWidgetRequest`, `WidgetListResponse`. Errors use `Problem`. Request objects set `additionalProperties: false` and `x-oapi-codegen-extra-tags` for `validate`.

## 2. Generate

```bash
make generate
```

## 3. Migration and queries

`migrations/NNNNN_widgets.sql` with goose Up and Down. Mirror the table in `sql/schema/schema.sql`. Add `sql/queries/widgets.sql` (`GetWidget :one`, `ListWidgets :many`, `CreateWidget :one`). List order includes a unique tie-break. Then `make generate` again.

## 4. Package `internal/widget/`

| File | Role |
|------|------|
| `types.go` | Domain struct |
| `errors.go` | `ErrNotFound` and other sentinels |
| `repository.go` | `Repository` and `SQLRepository` |
| `service.go` | Use cases. No sqlc types. |
| `mapper.go` | Domain ↔ OpenAPI |
| `handler.go` | Decode, validate, call service, `httpx` |

## 5. Wire

Construct the repository, service, and handler in `internal/app/app.go`. Delegate from `internal/transport/api.go`.

## 6. Test

Unit-test service rules with a fake repository. Add `integration/` coverage that skips without `DATABASE_URL` and uses `uuid.NewString()` in unique fields.

```bash
.cursor/skills/go-be-template/scripts/validate.sh
```
