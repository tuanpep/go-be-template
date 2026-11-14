package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/tuanbt/go-be-template/internal/app"
	"github.com/tuanbt/go-be-template/internal/config"
	"github.com/tuanbt/go-be-template/internal/platform/log"
	"github.com/tuanbt/go-be-template/internal/platform/migrate"
	"github.com/tuanbt/go-be-template/internal/platform/postgres"
	"github.com/tuanbt/go-be-template/migrations"
)

func TestUserAPI(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	if err := migrate.Up(ctx, pool, migrations.FS, "."); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	cfg := config.Config{
		HTTPAddr:    ":8080",
		DatabaseURL: databaseURL,
		LogLevel:    "error",
		LogFormat:   "text",
	}
	logger := log.New(cfg.LogLevel, cfg.LogFormat)
	slog.SetDefault(logger)

	handler, testPool, err := app.HandlerForTest(ctx, cfg, logger)
	if err != nil {
		t.Fatalf("app: %v", err)
	}
	defer testPool.Close()

	srv := httptest.NewServer(handler)
	defer srv.Close()

	t.Run("health", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/healthz")
		if err != nil {
			t.Fatalf("healthz: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("healthz status %d", resp.StatusCode)
		}
	})

	t.Run("ready", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/readyz")
		if err != nil {
			t.Fatalf("readyz: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("readyz status %d", resp.StatusCode)
		}
	})

	var userID string
	email := fmt.Sprintf("integration-%s@example.com", uuid.NewString())

	t.Run("create normalizes email", func(t *testing.T) {
		mixed := strings.ToUpper(email[:5]) + email[5:]
		createBody, _ := json.Marshal(map[string]string{
			"email": mixed,
			"name":  "Integration User",
		})
		resp, err := http.Post(srv.URL+"/v1/users", "application/json", bytes.NewReader(createBody))
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("create status %d body %s", resp.StatusCode, b)
		}
		var created map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got, _ := created["email"].(string); got != email {
			t.Fatalf("email: got %q want %q", got, email)
		}
		userID, _ = created["id"].(string)
		if userID == "" {
			t.Fatal("expected id")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"email": email, "name": "Other"})
		resp, err := http.Post(srv.URL+"/v1/users", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("duplicate status %d", resp.StatusCode)
		}
	})

	t.Run("get", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/v1/users/" + userID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("get status %d", resp.StatusCode)
		}
	})

	t.Run("get not found", func(t *testing.T) {
		id := uuid.New().String()
		resp, err := http.Get(srv.URL + "/v1/users/" + id)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("get status %d", resp.StatusCode)
		}
	})

	t.Run("bad json", func(t *testing.T) {
		resp, err := http.Post(srv.URL+"/v1/users", "application/json", strings.NewReader("{bad"))
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status %d", resp.StatusCode)
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/v1/users/not-a-uuid")
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status %d", resp.StatusCode)
		}
	})

	t.Run("name too long", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email": "long-" + uuid.NewString() + "@example.com",
			"name":  strings.Repeat("a", 201),
		})
		resp, err := http.Post(srv.URL+"/v1/users", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status %d", resp.StatusCode)
		}
	})

	t.Run("page size too large", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/v1/users?page=1&page_size=101")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status %d", resp.StatusCode)
		}
	})

	t.Run("update name", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"name": "Updated User"})
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/v1/users/"+userID, bytes.NewReader(body))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			t.Fatalf("update status %d body %s", resp.StatusCode, b)
		}
		var updated map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got, _ := updated["name"].(string); got != "Updated User" {
			t.Fatalf("name: got %q", got)
		}
	})

	t.Run("update duplicate email", func(t *testing.T) {
		other := fmt.Sprintf("other-%s@example.com", uuid.NewString())
		createBody, _ := json.Marshal(map[string]string{"email": other, "name": "Other"})
		created, err := http.Post(srv.URL+"/v1/users", "application/json", bytes.NewReader(createBody))
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		created.Body.Close()
		if created.StatusCode != http.StatusCreated {
			t.Fatalf("create status %d", created.StatusCode)
		}

		body, _ := json.Marshal(map[string]string{"email": other})
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/v1/users/"+userID, bytes.NewReader(body))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("update status %d", resp.StatusCode)
		}
	})

	t.Run("update missing", func(t *testing.T) {
		body := strings.NewReader(`{"name":"Missing"}`)
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/v1/users/"+uuid.NewString(), body)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("update: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("update status %d", resp.StatusCode)
		}
	})

	t.Run("delete", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/users/"+userID, nil)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("delete: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("delete status %d", resp.StatusCode)
		}
	})

	t.Run("delete missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/v1/users/"+userID, nil)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("delete: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("delete status %d", resp.StatusCode)
		}
	})

	t.Run("list", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/v1/users?page=1&page_size=10")
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("list status %d", resp.StatusCode)
		}
	})
}
