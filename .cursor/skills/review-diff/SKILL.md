---
name: review-diff
description: Reviews pull requests and local diffs against this template's architecture, OpenAPI, sqlc, and CI. Use when the user asks for a code review, PR review, or audit.
---

# Review a diff

| Label | Meaning |
|-------|---------|
| Critical | Secrets in the tree, SQL injection, migrate on API boot, broken data integrity |
| Major | Layer violation, spec drift, wrong HTTP status or error shape, missing integration test for a new API |
| Minor | Naming, docs, optional simplification |

Reject hand-edited `api.gen.go` or `internal/store/db`, SQL or `*db.Queries` in a service or handler, auto-migrate in `cmd/api`, a route missing from OpenAPI, and Problem JSON built outside `httpx`.

When the diff touches the spec or SQL, run `make generate` and `git diff --exit-code internal/store/db internal/transport/openapi`. Run `go test` for the packages in the diff. CI already runs the race detector, golangci-lint, and govulncheck.
