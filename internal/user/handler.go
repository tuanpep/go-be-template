package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/tuanbt/go-be-template/internal/platform/httpx"
	"github.com/tuanbt/go-be-template/internal/transport/openapi"
)

// Handler serves user-related HTTP operations.
type Handler struct {
	log      *slog.Logger
	svc      *Service
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, svc *Service) *Handler {
	return &Handler{
		log:      log,
		svc:      svc,
		validate: validator.New(),
	}
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, params openapi.ListUsersParams) {
	page, pageSize := 1, 20
	if params.Page != nil {
		page = *params.Page
	}
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	limit, offset, err := httpx.ParsePage(page, pageSize)
	if err != nil {
		httpx.BadRequest(w, "invalid page")
		return
	}

	rows, total, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		h.internal(w, r, "list users", err)
		return
	}

	items := make([]openapi.User, 0, len(rows))
	for _, row := range rows {
		u, err := ToAPI(row)
		if err != nil {
			h.internal(w, r, "map user", err)
			return
		}
		items = append(items, u)
	}

	httpx.WriteJSON(w, http.StatusOK, openapi.UserListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body openapi.CreateUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	if err := h.validate.Struct(body); err != nil {
		httpx.BadRequest(w, httpx.ValidationDetail(err))
		return
	}

	created, err := h.svc.Create(r.Context(), string(body.Email), body.Name)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			httpx.BadRequest(w, "invalid user input")
			return
		}
		if errors.Is(err, ErrDuplicateEmail) {
			httpx.Conflict(w, "email already registered")
			return
		}
		h.internal(w, r, "create user", err)
		return
	}

	resp, err := ToAPI(created)
	if err != nil {
		h.internal(w, r, "map user", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, userID openapi.UserId) {
	id, err := ParseID(userID)
	if err != nil {
		httpx.BadRequest(w, "invalid user id")
		return
	}

	found, err := h.svc.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			httpx.NotFound(w, "user not found")
			return
		}
		h.internal(w, r, "get user", err)
		return
	}

	resp, err := ToAPI(found)
	if err != nil {
		h.internal(w, r, "map user", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) internal(w http.ResponseWriter, r *http.Request, msg string, err error) {
	h.log.Error(msg,
		slog.String("error", err.Error()),
		slog.String("request_id", chimw.GetReqID(r.Context())),
	)
	httpx.InternalError(w)
}
