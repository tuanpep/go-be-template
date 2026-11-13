package user

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
