# Wiring

Only `internal/app` constructs the pool, `*db.Queries`, repositories, services, handlers, `health.Gate`, `transport.API`, the router, and `*http.Server`.

## Constructors

```
postgres.NewPool → db.New → NewSQLRepository → NewService → NewHandler → transport.NewAPI → NewRouter
```

`NewPool` does not ping. Pass interfaces into services. Pass concrete handlers into transport. No package-level wiring and no `init`. No uber/fx, wire, or dig unless the user asks.

## Process

1. `New` builds the pool and the server.
2. `ListenAndServe` runs in a goroutine.
3. SIGINT or SIGTERM closes the readiness gate, waits `ReadinessDrain` (2s), then `Shutdown` (`ShutdownTimeout`, 10s).
4. A second signal calls `Close` immediately.
5. `defer pool.Close()` runs after shutdown.

`cmd/api/main.go` loads config and calls `app.RunWithConfig`. It stays under 80 lines and does not register routes.

`cmd/migrate` pings, then runs goose, then exits.

Tests use `app.HandlerForTest`. The caller migrates and closes the pool.
