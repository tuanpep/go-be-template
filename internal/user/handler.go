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
	svc      Service
	validate *validator.Validate
}

func NewHandler(log *slog.Logger, svc Service) *Handler {
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
		items = append(items, ToAPI(row))
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

	created, err := h.svc.Create(r.Context(), CreateInput{
		Email: string(body.Email),
		Name:  body.Name,
	})
	if err != nil {
		h.writeUseCaseError(w, r, "create user", err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, ToAPI(created))
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, userID openapi.UserId) {
	id, err := ParseID(userID)
	if err != nil {
		httpx.BadRequest(w, "invalid user id")
		return
	}

	var body openapi.UpdateUserRequest
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

	in := UpdateInput{Name: body.Name}
	if body.Email != nil {
		v := string(*body.Email)
		in.Email = &v
	}
	updated, err := h.svc.Update(r.Context(), id, in)
	if err != nil {
		h.writeUseCaseError(w, r, "update user", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, ToAPI(updated))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, userID openapi.UserId) {
	id, err := ParseID(userID)
	if err != nil {
		httpx.BadRequest(w, "invalid user id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.writeUseCaseError(w, r, "delete user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, userID openapi.UserId) {
	id, err := ParseID(userID)
	if err != nil {
		httpx.BadRequest(w, "invalid user id")
		return
	}

	found, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.writeUseCaseError(w, r, "get user", err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, ToAPI(found))
}

func (h *Handler) writeUseCaseError(w http.ResponseWriter, r *http.Request, op string, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		httpx.BadRequest(w, "invalid user input")
	case errors.Is(err, ErrNotFound):
		httpx.NotFound(w, "user not found")
	case errors.Is(err, ErrDuplicateEmail):
		httpx.Conflict(w, "email already registered")
	default:
		h.internal(w, r, op, err)
	}
}

func (h *Handler) internal(w http.ResponseWriter, r *http.Request, msg string, err error) {
	h.log.Error(msg,
		slog.String("error", err.Error()),
		slog.String("request_id", chimw.GetReqID(r.Context())),
	)
	httpx.InternalError(w)
}
