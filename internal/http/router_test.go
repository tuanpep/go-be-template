package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/tuanbt/go-be-template/internal/transport"
	"github.com/tuanbt/go-be-template/internal/user"
)

type readyCheck struct{}

func (readyCheck) Ping(context.Context) error { return nil }

type fakeUsers struct{}

func (fakeUsers) Get(context.Context, uuid.UUID) (user.User, error) {
	return user.User{}, user.ErrNotFound
}
func (fakeUsers) List(context.Context, int32, int32) ([]user.User, int64, error) {
	return nil, 0, nil
}
func (fakeUsers) Create(context.Context, string, string) (user.User, error) {
	return user.User{}, nil
}

func TestRouter_InvalidUserIDIsProblemJSON(t *testing.T) {
	t.Parallel()
	log := slog.New(slog.DiscardHandler)
	api := transport.NewAPI(log, readyCheck{}, user.NewHandler(log, user.NewService(fakeUsers{})))
	router := NewRouter(log, api)

	req := httptest.NewRequest(http.MethodGet, "/v1/users/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type %q body %s", ct, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"status":400`) {
		t.Fatalf("body %s", rec.Body.String())
	}
}

func TestRouter_PageSizeAboveMaxIsBadRequest(t *testing.T) {
	t.Parallel()
	log := slog.New(slog.DiscardHandler)
	api := transport.NewAPI(log, readyCheck{}, user.NewHandler(log, user.NewService(fakeUsers{})))
	router := NewRouter(log, api)

	req := httptest.NewRequest(http.MethodGet, "/v1/users?page=1&page_size=101", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}
