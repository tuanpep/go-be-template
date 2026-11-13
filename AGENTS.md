# Agent guide

Skills live in `.cursor/skills/`. Read **go-be-template** first. It routes to one doc and one verify script. Layer skills are thin procedures with `paths`, so they surface only for matching files.

| Skill | When |
|-------|------|
| **go-be-template** | Any change in this repository |
| go-be-template-structure | Packages, imports, new feature folders |
| go-be-template-openapi | `api/openapi.yaml`, codegen, routes |
| go-be-template-data | Migrations, sqlc, repositories |
| go-be-template-http | Chi, middleware, handlers, httpx |
| go-be-template-wiring | `internal/app`, `cmd/api`, `cmd/migrate` |
| go-be-template-errors | Sentinels, slog, Problem responses |
| go-be-template-testing | Tests, CI, golangci |
| go-be-template-operations | Docker, env, health, migrate job |
| go-be-template-security | Secrets, SQL safety, public routes |
| go-be-template-review | Code review, PR review, audit |

**Before done:** `.cursor/skills/go-be-template/scripts/validate.sh`

**One layer:** `.cursor/skills/go-be-template/scripts/verify-<section>.sh`

**New feature:** `.cursor/skills/go-be-template/examples.md`

Cursor attaches `.cursor/rules/go-be-template.mdc` when Go, SQL, OpenAPI, or deploy files are in context.
