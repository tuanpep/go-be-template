# Agent guide

Each skill is `.cursor/skills/<name>/SKILL.md`. The rule with the same name attaches on the same files and points at that skill.

| Skill and rule | Action |
|----------------|--------|
| **scaffold-feature** | Add a feature or endpoint set |
| place-packages | Create a package, move files, or fix an import cycle |
| edit-openapi | Edit `api/openapi.yaml`, codegen, or routes |
| add-persistence | Add a migration, sqlc query, or repository |
| write-http | Edit Chi, middleware, handlers, or httpx |
| wire-app | Edit `internal/app`, `cmd/api`, or `cmd/migrate` |
| map-errors | Edit sentinels, slog, or Problem responses |
| write-tests | Add tests, CI, or golangci config |
| configure-runtime | Edit Docker, env, health, or the migrate job |
| guard-security | Touch secrets, SQL safety, or public routes |
| review-diff | Review a diff, pull request, or audit |

`scaffold-feature` and `review-diff` apply when the request matches. The other rules apply when their globs match.

After a spec or SQL change: `make generate`, then `git diff --exit-code internal/store/db internal/transport/openapi`. Run `go test` for the packages you changed.
