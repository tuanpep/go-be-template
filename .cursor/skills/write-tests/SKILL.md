---
name: write-tests
description: Writes unit tests, integration tests, and CI checks. Use when adding tests or editing .github/workflows, .golangci.yml, or Makefile test targets.
paths:
  - "**/*_test.go"
  - "integration/**"
  - ".github/**"
  - ".golangci.yml"
---

# Write tests

| Test | Where | Rule |
|------|-------|------|
| Unit | Next to the package | Fake `Repository`. `t.Parallel()` is fine. |
| Integration | `integration/` | Skip when `DATABASE_URL` is unset. |
| Generated | — | Leave generated packages untested. |

Integration setup is `postgres.NewPool`, `migrate.Up`, `app.HandlerForTest`, and `httptest`. Close pools with `defer`. Use `uuid.NewString()` in emails and ids. `t.Name()` stays the same across runs and collides on a shared database.

CI order: setup-go from `go.mod` with `check-latest: true`, `go tool` sqlc and oapi-codegen then `git diff --exit-code` on the generated trees, golangci-lint v2.9.0, `go tool govulncheck ./...`, `go test -race ./...`, then build `cmd/api` and `cmd/migrate`. Tool versions are the `tool` directives in `go.mod`. After a spec or SQL edit, run `make generate` before tests.

Run `go test` for the packages you changed.
