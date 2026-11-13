# Code review

## Severity

| Label | Meaning |
|-------|---------|
| Critical | Secrets in the tree, SQL injection, migrate on API boot, broken data integrity |
| Major | Layer violation, spec drift, wrong HTTP status or error shape, missing integration test for a new API |
| Minor | Naming, docs, optional simplification |

## Reject

1. Hand-edited `api.gen.go` or `internal/store/db/*`.
2. SQL or `*db.Queries` in a service or handler.
3. Auto-migrate in `cmd/api`.
4. Route that is not in OpenAPI.
5. Problem JSON built without `httpx`.

## Also run

| Diff touches | Script |
|--------------|--------|
| `openapi.yaml` | `verify-openapi.sh` |
| `migrations/`, `sql/` | `verify-data.sh` |
| `internal/app` | `verify-wiring.sh` |
| router or `httpx` | `verify-http.sh` |

Suggest `scripts/validate.sh` before merge.
