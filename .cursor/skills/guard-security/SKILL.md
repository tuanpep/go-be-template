---
name: guard-security
description: Guards secrets, sqlc-only SQL, and public routes. Use when adding auth, env secrets, raw SQL, outbound HTTP, or PII fields.
paths:
  - "api/openapi.yaml"
  - ".env.example"
  - "internal/config/**/*.go"
  - "internal/http/**/*.go"
  - "internal/**/repository.go"
---

# Guard security

This template has no authentication. The `/v1/users` routes are public samples. Say that when adding a public route.

SQL is sqlc placeholders. A dynamic `ORDER BY` uses a column allowlist in the repository. `DATABASE_URL` is env-only and required. Leave it out of logs, OpenAPI examples, and Problem `detail`. Validate bodies with the generated `validate` tags and reject unknown JSON fields.

CORS is off. A browser client needs an allowlist. Rate limiting is outside the process. Add it at the proxy, or as Chi middleware in `internal/http`, before exposing the sample on the internet. An outbound call to a user-supplied URL needs an allowlist.

When auth is requested, parse JWT and store claims in `internal/platform/auth`, enforce them in `internal/http` middleware, and read the user id from context in the handler. Leave JWTs out of sqlc models. Hash passwords with `golang.org/x/crypto` and compare secrets in constant time.

Flag for review: authz, uploads, raw SQL, user-controlled URLs, and PII in logs.
