---
name: configure-runtime
description: Configures Docker, the migrate job, env, and health probes. Use when editing the Dockerfile, Makefile, compose file, .env.example, config, health, or cmd/migrate.
paths:
  - "Dockerfile"
  - "docker-compose.yml"
  - "Makefile"
  - ".env.example"
  - "cmd/migrate/**"
  - "internal/config/**/*.go"
  - "internal/platform/health/**/*.go"
  - "internal/platform/postgres/**/*.go"
---

# Configure runtime

`DATABASE_URL` is required. The process exits when it is unset. `make run` and `make migrate` load `.env` when that file exists, and otherwise pass the compose URL. The binary does not read `.env`. Defaults are `HTTP_ADDR=:8080`, `LOG_LEVEL=info`, and `LOG_FORMAT=json`. Document a new variable in `.env.example`. `.env` stays gitignored.

| Binary | Migrations | Database |
|--------|------------|----------|
| `cmd/migrate` | goose up | Ping, then exit on failure |
| `cmd/api` | None | Pool with no startup ping |

Deploy migrate, then start the API. `/healthz` returns 200 while the process is listening. `/readyz` returns 200 when the pool answers and the process is accepting traffic, otherwise 503.

On SIGTERM the API closes the readiness gate, waits 2s, then shuts down for up to 10s. A second signal closes the listener.

The image is multi-stage, `CGO_ENABLED=0`, user `appuser` (uid 10001). `.dockerignore` excludes `.env`. Migrations in the image are for visibility. The migrate binary uses the embedded FS.
