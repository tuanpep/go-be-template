package user

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/tuanbt/go-be-template/internal/transport/openapi"
)

const stubID = "11111111-1111-1111-1111-111111111111"

func mustUserID(t *testing.T) openapi.UserId {
	t.Helper()
	id, err := uuid.Parse(stubID)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestCreateUser_RejectsOverlongName(t *testing.T) {
	t.Parallel()
	h := NewHandler(slog.New(slog.DiscardHandler), NewService(&stubRepo{}))
	body := `{"email":"ada@example.com","name":"` + strings.Repeat("a", 201) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_RejectsEmptyBody(t *testing.T) {
	t.Parallel()
	h := NewHandler(slog.New(slog.DiscardHandler), NewService(&stubRepo{}))
	req := httptest.NewRequest(http.MethodPatch, "/v1/users/"+stubID, strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req, mustUserID(t))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateUser_RejectsUnknownField(t *testing.T) {
	t.Parallel()
	h := NewHandler(slog.New(slog.DiscardHandler), NewService(&stubRepo{}))
	req := httptest.NewRequest(http.MethodPatch, "/v1/users/"+stubID, strings.NewReader(`{"name":"Ada","role":"admin"}`))
	rec := httptest.NewRecorder()
	h.UpdateUser(rec, req, mustUserID(t))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}

func TestCreateUser_RejectsUnknownField(t *testing.T) {
	t.Parallel()
	h := NewHandler(slog.New(slog.DiscardHandler), NewService(&stubRepo{}))
	body := `{"email":"ada@example.com","name":"Ada","role":"admin"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/users", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.CreateUser(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status %d body %s", rec.Code, rec.Body.String())
	}
}
