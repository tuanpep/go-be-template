# Operations

## Config

All deploy-varying config comes from the environment.

- `DATABASE_URL` is required. The process exits when it is unset. There is no compiled-in credential.
- `make run` and `make migrate` load `.env` when it exists, otherwise they pass the compose URL. The binary does not read `.env`.
- Defaults: `HTTP_ADDR=:8080`, `LOG_LEVEL=info`, `LOG_FORMAT=json`.

Document new variables in `.env.example`. `.env` stays gitignored.

## Processes

| Binary | Migrations | Database |
|--------|------------|----------|
| `cmd/migrate` | goose up | Ping, then exit on failure |
| `cmd/api` | No | Pool without a startup ping |

Deploy migrate, then start the API.

## Probes

| Probe | Path | Result |
|-------|------|--------|
| Liveness | `/healthz` | 200 while the process is listening |
| Readiness | `/readyz` | 200 when the pool answers and the process is not shutting down, else 503 |

On SIGTERM the API closes the readiness gate, waits 2s, then shuts down for up to 10s. A second signal closes the listener.

## Image

Multi-stage Dockerfile, `CGO_ENABLED=0`, user `appuser` (uid 10001). `.dockerignore` excludes `.env`. The runtime image copies the binaries. Migrations in the image are for visibility. The migrate binary uses the embedded FS.
