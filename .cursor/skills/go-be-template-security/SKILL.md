---
name: go-be-template-security
description: Applies secrets, sqlc-only SQL, and public-route rules for go-be-template. Use when adding auth, env secrets, raw SQL, outbound HTTP, or PII fields.
paths:
  - "internal/**/*.go"
  - "api/openapi.yaml"
  - ".env.example"
---

# Security

1. Read [security.md](../go-be-template/docs/security.md).
2. Parameterized sqlc only. Say when a new route is public. The template has no auth until one is requested.
3. Run `.cursor/skills/go-be-template/scripts/verify-security.sh`.
