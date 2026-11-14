---
name: grow-feature
description: >-
  Guides how a feature grows across this repo (OpenAPI, SQL, packages, tests,
  runtime, optional second process). Use when a feature is large, spans multiple
  packages, or needs stage-based scaling beyond scaffold-feature.
paths:
  - "cmd/**"
  - "Dockerfile"
  - "docker-compose.yml"
  - ".github/**"
  - "internal/*/*/**"
---

# Grow a feature

One public `Service` per feature. Layout rules: **codebase-architecture**. OSS stage examples: [codebase-architecture/references/oss-repos.md](../codebase-architecture/references/oss-repos.md).

New feature starts with **scaffold-feature** (`internal/user`). Apply the first matching stage only.

## Instructions

| Stage | Signal | Change |
|-------|--------|--------|
| 1 | More use cases, same noun | OpenAPI + sql + `Service` methods; same package |
| 2 | `service.go` too large | Split files (`create.go`, `list.go`); same `Service` |
| 3 | Own lifecycle (scheduler, strategy, state machine) | `internal/<feature>/<behavior>/`; parent exposes small interface |
| 4 | New noun | New `internal/<feature>/`; full **scaffold-feature** |
| 5 | Cache, audit, metrics on a port | Decorator with same interface; generic bits in `platform` |
| 6 | Own scale/DB/stable contract | New `cmd/<name>` + **configure-runtime** + CI; last resort |

Stage 3+: sqlc stays `internal/store/db`; `transport` still delegates; background work wired in `app` and stopped on shutdown.

Per-layer work: **edit-openapi**, **add-persistence**, **write-http**, **wire-app**, **map-errors**, **write-tests**, **guard-security** (not repeated here).