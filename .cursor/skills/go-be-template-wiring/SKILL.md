---
name: go-be-template-wiring
description: Guides composition root and process lifecycle for go-be-template. Use when editing internal/app, cmd/api, cmd/migrate, or how dependencies are constructed.
paths:
  - "internal/app/**/*.go"
  - "cmd/**/*.go"
---

# Wiring

1. Read [application-wiring.md](../go-be-template/docs/application-wiring.md).
2. Construct pool → sqlc → repository → service → handler → transport → router only in `internal/app`.
3. Run `.cursor/skills/go-be-template/scripts/verify-wiring.sh`.
