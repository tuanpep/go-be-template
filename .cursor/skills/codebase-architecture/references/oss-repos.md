# OSS layouts (GitHub)

Verified via GitHub API / in-repo docs. Maps to this template's skills, not copy-paste targets.

## At a glance

| Repo | Feature unit | This template |
|------|----------------|---------------|
| [grafana/grafana](https://github.com/grafana/grafana) | `pkg/services/<domain>/`, e.g. [ngalert](https://github.com/grafana/grafana/tree/main/pkg/services/ngalert) subdirs | `internal/<feature>/` + stage 3 subpackages |
| [ory/kratos](https://github.com/ory/kratos) | [selfservice/flow/*](https://github.com/ory/kratos/tree/master/selfservice/flow), [strategy/*](https://github.com/ory/kratos/tree/master/selfservice/strategy) | Same package + `/<behavior>` for strategies |
| [mattermost/mattermost](https://github.com/mattermost/mattermost) | [api4 → app → store](https://github.com/mattermost/mattermost/tree/master/server/channels) | We use feature repo, not global `store` |
| [go-gitea/gitea](https://github.com/go-gitea/gitea) | [routers](https://github.com/go-gitea/gitea/tree/main/routers) + [services](https://github.com/go-gitea/gitea/tree/main/services) + [models](https://github.com/go-gitea/gitea/tree/main/models) | `transport` + colocated feature, not a 3-way split |
| [zitadel/zitadel](https://github.com/zitadel/zitadel) | [internal/command, query, api](https://github.com/zitadel/zitadel/tree/main/internal); [backend/v3](https://github.com/zitadel/zitadel/tree/main/backend/v3) | Not default CQRS |
| [bytebase/bytebase](https://github.com/bytebase/bytebase) | [api/v1/*_service.go](https://github.com/bytebase/bytebase/tree/main/backend/api/v1) + [store](https://github.com/bytebase/bytebase/tree/main/backend/store) | OpenAPI + thin transport; fat services not copied |
| [temporalio/temporal](https://github.com/temporalio/temporal) | [service/*](https://github.com/temporalio/temporal/tree/main/service) processes | **grow-feature** stage 6 only |

## Docs worth reading in-tree

- Grafana [contribute/backend/services.md](https://github.com/grafana/grafana/blob/main/contribute/backend/services.md) — service + Wire.
- Mattermost [app/doc.go](https://github.com/mattermost/mattermost/blob/master/server/channels/app/doc.go), [store/doc.go](https://github.com/mattermost/mattermost/blob/master/server/channels/store/doc.go).
- Zitadel [backend/v3/doc.go](https://github.com/zitadel/zitadel/blob/main/backend/v3/doc.go).
- Temporal [history-service.md](https://github.com/temporalio/temporal/blob/main/docs/architecture/history-service.md).

## grow-feature stages ↔ OSS

| Stage | OSS example |
|-------|----------------|
| 3 | Grafana ngalert `schedule`/`state`; Kratos `strategy/*`; Temporal `history/workflow` |
| 5 | Mattermost store decorator layers (here: per-repository decorator only) |
| 6 | Temporal `service/frontend|history|matching|worker` |