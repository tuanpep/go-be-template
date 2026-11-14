---
name: wire-app
description: Wires the composition root and process lifecycle. Use when editing internal/app, cmd/api, cmd/migrate, or how dependencies are constructed.
paths:
  - "internal/app/**/*.go"
  - "cmd/**/*.go"
---

# Wire the app

Only `internal/app` constructs the pool, `*db.Queries`, repositories, services, handlers, `health.Gate`, `transport.API`, the router, and `*http.Server`.

```
postgres.NewPool → db.New → NewSQLRepository(pool, queries) → NewService → NewHandler → transport.NewAPI → NewRouter
```

`NewPool` does not ping. Pass interfaces into services and concrete handlers into transport. Keep wiring out of `init` and package globals.

On SIGINT or SIGTERM, close the readiness gate, wait `ReadinessDrain` (2s), then `Shutdown` for `ShutdownTimeout` (10s). A second signal calls `Close`. `defer pool.Close()` runs after shutdown.

`cmd/api/main.go` loads config and calls `app.RunWithConfig`. It stays under 80 lines and does not register routes. `cmd/migrate` pings, runs goose, and exits. Tests use `app.HandlerForTest`; the caller migrates and closes the pool.
