---
name: go-be-template-testing
description: Defines unit tests, integration tests, and CI for go-be-template. Use when adding tests or editing .github/workflows, .golangci.yml, or Makefile test targets.
paths:
  - "**/*_test.go"
  - "integration/**"
  - ".github/**"
  - ".golangci.yml"
---

# Testing

1. Read [testing-ci.md](../go-be-template/docs/testing-ci.md).
2. Unit-test domain rules with a fake repository. Integration tests skip when `DATABASE_URL` is unset and use a fresh id each run.
3. Run `.cursor/skills/go-be-template/scripts/verify-testing.sh`.
