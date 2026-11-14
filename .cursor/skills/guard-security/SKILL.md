---
name: guard-security
description: >-
  Applies secrets hygiene, sqlc-only SQL, and public-route conventions for this
  template. Use when adding auth, editing repositories, HTTP middleware,
  OpenAPI, or .env.example.
paths:
  - "api/openapi.yaml"
  - ".env.example"
  - "internal/config/**"
  - "internal/http/**"
  - "internal/**/repository.go"
---

# Guard security

No auth in template; `/v1/users` is a public sample.

## Instructions

- No `DATABASE_URL` or secrets in logs or Problem bodies.
- CORS off; rate limit at proxy or Chi if public. User-supplied outbound URLs → allowlist.

**add-persistence** / repositories: sqlc-only SQL; dynamic `ORDER BY` → allowlist.

**write-http** / handlers: reject unknown JSON fields (`DisallowUnknownFields`).

Auth (when asked): new `internal/platform/auth`, middleware in `http`, claims on `context`; hash secrets with `x/crypto`, constant-time compare.