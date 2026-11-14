# Skills

Project skills follow the [Agent Skills](https://agentskills.io/specification) format: each skill is a directory named for its `name` field with a required `SKILL.md` (YAML frontmatter + markdown body). Detailed reference material lives one level deep under `references/`, `scripts/`, or `assets/` when needed.

Cursor adds an optional [`paths`](https://cursor.com/docs/skills#scoping-a-skill-to-specific-files) frontmatter field (glob scoping). That field is not in the core agentskills.io schema; `skills-ref validate` reports it unless you strip `paths` for portable packages.

| Skill | When |
|-------|------|
| **codebase-architecture** | Layout, imports, package moves ([references/oss-repos.md](codebase-architecture/references/oss-repos.md)) |
| **scaffold-feature** | New resource |
| **grow-feature** | Feature outgrows one package |
| **edit-openapi** | `api/openapi.yaml`, codegen |
| **add-persistence** | Migrations, sqlc, repository |
| **write-http** | Router, httpx, handlers |
| **wire-app** | `internal/app`, `cmd/*` |
| **map-errors** | Sentinels, Problem JSON |
| **write-tests** | Tests, CI |
| **configure-runtime** | Docker, env, health, migrate job |
| **guard-security** | Secrets, SQL, auth |
| **review-diff** | PR / diff review (invoke with `/review-diff`; not auto on every Go file) |

After spec or SQL change: `make generate`, `git diff --exit-code internal/store/db internal/transport/openapi`, `go test` (see [AGENTS.md](../../AGENTS.md)).