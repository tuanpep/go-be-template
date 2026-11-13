# Security

The template ships without authentication. Say that when adding a public route. `POST /v1/users` and `GET /v1/users` are public samples.

## Rules

- SQL is sqlc placeholders. No `fmt.Sprintf` SQL. Dynamic `ORDER BY` uses a column allowlist in the repository.
- `DATABASE_URL` is env-only and required. Do not log it. Do not put secrets in OpenAPI examples or Problem `detail`.
- Validate bodies with the generated `validate` tags. Reject unknown JSON fields.
- CORS is off. An allowlist is required for browser clients. Never `*` with credentials.
- Rate limiting is not in the process. Add it at the proxy, or as Chi middleware in `internal/http`, before exposing the sample on the internet.
- Outbound HTTP to a user-supplied URL needs an allowlist. There is no such call today.

## When auth is requested

```
internal/platform/auth/      JWT parse, context claims
internal/http/middleware.go  RequireAuth
feature handler              user id from context
```

Do not store JWTs in sqlc models. Use `golang.org/x/crypto` for password hashes. Compare secrets in constant time.

Flag for review: authz, uploads, raw SQL, user-controlled URLs, PII in logs.
