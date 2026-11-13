package transport

import (
	"errors"
	"log/slog"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/tuanbt/go-be-template/internal/platform/health"
	"github.com/tuanbt/go-be-template/internal/platform/httpx"
	"github.com/tuanbt/go-be-template/internal/transport/openapi"
	"github.com/tuanbt/go-be-template/internal/user"
)

// API implements the generated OpenAPI server by delegating to domain handlers.
type API struct {
	log   *slog.Logger
	check health.Checker
	users *user.Handler
}

func NewAPI(log *slog.Logger, check health.Checker, users *user.Handler) *API {
	return &API{log: log, check: check, users: users}
}

var _ openapi.ServerInterface = (*API)(nil)

func (a *API) GetHealth(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, openapi.HealthResponse{Status: "ok"})
}

func (a *API) GetReadiness(w http.ResponseWriter, r *http.Request) {
	if err := a.check.Ping(r.Context()); err != nil {
		detail := "database unreachable"
		if errors.Is(err, health.ErrUnavailable) {
			detail = "not ready"
		} else {
			a.log.Warn("readiness check failed",
				slog.String("error", err.Error()),
				slog.String("request_id", chimw.GetReqID(r.Context())),
			)
		}
		httpx.WriteProblem(w, http.StatusServiceUnavailable, "Service Unavailable", detail)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, openapi.HealthResponse{Status: "ok"})
}

func (a *API) ListUsers(w http.ResponseWriter, r *http.Request, params openapi.ListUsersParams) {
	a.users.ListUsers(w, r, params)
}

func (a *API) CreateUser(w http.ResponseWriter, r *http.Request) {
	a.users.CreateUser(w, r)
}

func (a *API) GetUser(w http.ResponseWriter, r *http.Request, userID openapi.UserId) {
	a.users.GetUser(w, r, userID)
}
