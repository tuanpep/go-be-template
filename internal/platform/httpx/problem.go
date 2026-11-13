package httpx

import (
	"net/http"

	"github.com/tuanbt/go-be-template/internal/transport/openapi"
)

func WriteProblem(w http.ResponseWriter, status int, title, detail string) {
	p := openapi.Problem{Title: title, Status: status}
	if detail != "" {
		p.Detail = &detail
	}
	WriteJSON(w, status, p)
}

func BadRequest(w http.ResponseWriter, detail string) {
	WriteProblem(w, http.StatusBadRequest, "Bad Request", detail)
}

func NotFound(w http.ResponseWriter, detail string) {
	WriteProblem(w, http.StatusNotFound, "Not Found", detail)
}

func Conflict(w http.ResponseWriter, detail string) {
	WriteProblem(w, http.StatusConflict, "Conflict", detail)
}

func InternalError(w http.ResponseWriter) {
	WriteProblem(w, http.StatusInternalServerError, "Internal Server Error", "")
}
