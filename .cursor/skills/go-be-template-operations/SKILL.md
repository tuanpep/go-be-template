---
name: go-be-template-operations
description: Covers Docker, the migrate job, env config, and health probes for go-be-template. Use when editing the Dockerfile, Makefile, compose file, .env.example, or cmd/migrate.
paths:
  - "Dockerfile"
  - "docker-compose.yml"
  - "Makefile"
  - ".env.example"
  - "cmd/migrate/**"
  - "internal/config/**/*.go"
  - "internal/platform/health/**/*.go"
  - "internal/platform/postgres/**/*.go"
---

# Operations

1. Read [operations.md](../go-be-template/docs/operations.md).
2. New env vars go in `.env.example`. The API process does not migrate and does not invent `DATABASE_URL`.
3. Run `.cursor/skills/go-be-template/scripts/verify-operations.sh`.
