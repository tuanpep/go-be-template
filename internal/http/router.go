package httpserver

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/tuanbt/go-be-template/internal/platform/httpx"
	"github.com/tuanbt/go-be-template/internal/transport/openapi"
)

// NewRouter mounts API routes and standard middleware.
func NewRouter(log *slog.Logger, api openapi.ServerInterface) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Recoverer)
	r.Use(requestLog(log))
	r.Use(chimw.Timeout(RequestTimeout))
	r.Use(limitRequestBody(MaxRequestBodyBytes))

	openapi.HandlerWithOptions(api, openapi.ChiServerOptions{
		BaseRouter:       r,
		ErrorHandlerFunc: parameterError,
	})
	return r
}

// parameterError maps generated binding failures onto the Problem shape.
func parameterError(w http.ResponseWriter, _ *http.Request, err error) {
	var invalid *openapi.InvalidParamFormatError
	if errors.As(err, &invalid) {
		httpx.BadRequest(w, "invalid "+invalid.ParamName)
		return
	}
	var missing *openapi.RequiredParamError
	if errors.As(err, &missing) {
		httpx.BadRequest(w, missing.ParamName+" is required")
		return
	}
	httpx.BadRequest(w, "invalid parameter")
}

func limitRequestBody(max int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil && r.ContentLength != 0 {
				r.Body = http.MaxBytesReader(w, r.Body, max)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func requestLog(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", time.Since(start)),
				slog.String("request_id", chimw.GetReqID(r.Context())),
			)
		})
	}
}
