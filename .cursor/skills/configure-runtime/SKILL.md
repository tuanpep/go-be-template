---
name: configure-runtime
description: >-
  Documents Docker images, environment variables, health probes, and the migrate
  job for this template. Use when changing Dockerfile, docker-compose, Makefile,
  .env.example, config, or health gates.
paths:
  - "Dockerfile"
  - "docker-compose.yml"
  - "Makefile"
  - ".env.example"
  - "cmd/migrate/**"
  - "internal/config/**"
  - "internal/platform/health/**"
  - "internal/platform/postgres/**"
---

# Configure runtime

`DATABASE_URL` required. Defaults: `HTTP_ADDR=:8080`, `LOG_LEVEL=info`, `LOG_FORMAT=json`. Document env in `.env.example`; binaries do not read `.env` (Makefile may).

## Instructions

| Binary | DB | Migrations |
|--------|-----|------------|
| `cmd/migrate` | Ping | goose up, exit |
| `cmd/api` | Pool (no startup ping) | None |

Deploy migrate before API. `/healthz` = listening; `/readyz` = DB up (503 otherwise). Graceful shutdown sequence: **wire-app**.

Image: multi-stage, `CGO_ENABLED=0`, user `appuser` 10001.