---
name: go-be-template-errors
description: Defines sentinel errors, slog, and Problem responses for go-be-template. Use when editing errors.go, handler error branches, or platform/log.
paths:
  - "internal/**/errors.go"
  - "internal/**/handler.go"
  - "internal/platform/log/**/*.go"
  - "internal/platform/httpx/**/*.go"
---

# Errors

1. Read [errors-logging.md](../go-be-template/docs/errors-logging.md).
2. Map sentinels with `errors.Is` to `httpx`. Log unexpected errors with `request_id`.
3. Run `.cursor/skills/go-be-template/scripts/verify-errors.sh`.
