# Testing and CI

| Test | Where | Rule |
|------|-------|------|
| Unit | next to the package | Fake `Repository`. `t.Parallel()` is fine. |
| Integration | `integration/` | Skip when `DATABASE_URL` is unset. |
| Generated | — | Do not test generated packages directly. |

Integration setup: `postgres.NewPool`, `migrate.Up`, `app.HandlerForTest`, `httptest`. Close pools with `defer`.

Use `uuid.NewString()` in emails and ids. `t.Name()` is stable, so a second run against the same database collides.

## CI order

1. Checkout and setup-go from `go.mod`, `check-latest: true`.
2. `go tool sqlc` and `go tool oapi-codegen`, then `git diff --exit-code` on generated trees.
3. golangci-lint v2.9.0.
4. `go tool govulncheck ./...`.
5. `go test -race ./...` (includes integration when `DATABASE_URL` is set).
6. Build `cmd/api` and `cmd/migrate`.

Do not weaken the drift check. After a spec or SQL edit, run `make generate` before tests.

Tool versions are the `tool` directives in `go.mod`, not `@latest`.
