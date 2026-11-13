---
name: go-be-template-data
description: Applies goose migrations, sqlc, and pgx repositories for go-be-template. Use when writing migrations, sql/schema, sql/queries, or repository code.
paths:
  - "migrations/**"
  - "sql/**"
  - "sqlc.yaml"
  - "internal/**/repository.go"
  - "internal/store/**"
---

# Data

1. Read [database.md](../go-be-template/docs/database.md).
2. Migration, then `sql/schema`, then `sql/queries`, then `make generate`, then the repository.
3. Run `.cursor/skills/go-be-template/scripts/verify-data.sh`.
