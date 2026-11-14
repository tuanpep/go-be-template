---
name: review-diff
description: >-
  Reviews diffs against this template's architecture, OpenAPI, sqlc, security,
  and CI expectations. Use when reviewing pull requests, local diffs, or before
  merge.
disable-model-invocation: true
paths:
  - "api/**"
  - "sql/**"
  - "migrations/**"
  - ".github/**"
---

# Review a diff

Use **codebase-architecture** (imports, do-not list, request flow) and **grow-feature** stage order.

## Instructions

| Label | Examples |
|-------|----------|
| Critical | Secrets committed, SQL injection, migrate in `cmd/api`, broken integrity |
| Major | Violations of **codebase-architecture** § Do not, spec drift, wrong errors, new API without integration test |
| Minor | Naming, docs |

Extra reject checks (not spelled out in architecture):

- Hand-edited `internal/store/db` or `api.gen.go`
- Problem JSON outside `httpx`
- Cross-feature repository imports

Spec/SQL touched → `make generate` and clean generated diff (`git diff --exit-code internal/store/db internal/transport/openapi`). Run `go test` on affected packages.