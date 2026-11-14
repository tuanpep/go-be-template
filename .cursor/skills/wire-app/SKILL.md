---
name: wire-app
description: >-
  Wires the composition root and process lifecycle for this Go backend. Use when
  editing internal/app, cmd/api, cmd/migrate, or how dependencies are constructed.
paths:
  - "internal/app/**/*.go"
  - "cmd/**/*.go"
---

# Wire the app

Only `internal/app` constructs the pool, `*db.Queries`, repositories, services, handlers, `health.Gate`, `transport.API`, the router, and `*http.Server`.

## Instructions

```
postgres.NewPool → db.New → NewSQLRepository(pool, queries) → NewService(repo) (returns feature.Service) → NewHandler(log, svc) → transport.NewAPI → NewRouter
```

`NewPool` does not ping. Handlers take the feature `Service` interface, not the concrete struct. Keep wiring out of `init` and package globals. One `transport.API` field per feature handler; it only delegates.

On SIGINT or SIGTERM, close the readiness gate, wait `ReadinessDrain` (2s), then `Shutdown` for `ShutdownTimeout` (10s). A second signal calls `Close`. `defer pool.Close()` runs after shutdown.

`cmd/api/main.go` loads config and calls `app.RunWithConfig`. It stays under 80 lines and does not register routes. Tests use `app.HandlerForTest`; the caller migrates and closes the pool.

Migrate binary behavior, deploy order, and env: **configure-runtime** (do not run goose from `cmd/api`).