# Agent guide

Skills: `.cursor/skills/<name>/SKILL.md`. Index: [.cursor/skills/README.md](.cursor/skills/README.md).

| Skill | Action |
|-------|--------|
| **codebase-architecture** | Layout, imports, OSS map (`codebase-architecture/references/oss-repos.md`) |
| **scaffold-feature** | New feature or endpoint set |
| **grow-feature** | Scale a feature across the repo |
| edit-openapi | OpenAPI, codegen, routes |
| add-persistence | Migrations, sqlc, repository |
| write-http | Router, httpx, handlers |
| wire-app | `internal/app`, `cmd/*` |
| map-errors | Sentinels, Problem responses |
| write-tests | Tests, CI |
| configure-runtime | Docker, env, health, migrate |
| guard-security | Secrets, SQL safety, public routes |
| **review-diff** | PR or diff review (`/review-diff`) |

After spec or SQL change: `make generate`, then `git diff --exit-code internal/store/db internal/transport/openapi`. Run `go test` for changed packages.