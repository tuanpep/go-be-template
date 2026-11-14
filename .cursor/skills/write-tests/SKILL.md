---
name: write-tests
description: >-
  Covers unit tests, integration tests, and CI checks for this Go backend
  template. Use when adding *_test.go files, integration/ tests, or .github
  workflows.
paths:
  - "**/*_test.go"
  - "integration/**"
  - ".github/**"
---

# Write tests

## Instructions

| Kind | Rule |
|------|------|
| Unit | Fake `Repository`; `t.Parallel()` ok |
| Integration | `integration/`; skip without `DATABASE_URL`; `HandlerForTest` + migrate; unique fields via `uuid.NewString()` |
| Generated | Do not test sqlc/openapi packages |

CI: generate check, golangci-lint v2.9.0, govulncheck, `go test -race ./...`, build `cmd/api` and `cmd/migrate`.