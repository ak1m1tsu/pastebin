package response

import (
	"net/http"

	"github.com/go-chi/render"
)

func response(w http.ResponseWriter, r *http.Request, status int, v map[string]any) {
	render.Status(r, status)
	render.JSON(w, r, v)
}

func OK(w http.ResponseWriter, r *http.Request, v map[string]any) {
	response(w, r, http.StatusOK, v)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "method not allowed"}
	response(w, r, http.StatusMethodNotAllowed, v)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "resource not found"}
	response(w, r, http.StatusNotFound, v)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "internal server error"}
	response(w, r, http.StatusInternalServerError, v)
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "bad request"}
	response(w, r, http.StatusBadRequest, v)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "unauthorized"}
	response(w, r, http.StatusUnauthorized, v)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "forbidden"}
	response(w, r, http.StatusForbidden, v)
}

func Conflict(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "conflict"}
	response(w, r, http.StatusConflict, v)
}

func UnprocessableEntity(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "unprocessable entity"}
	response(w, r, http.StatusUnprocessableEntity, v)
}

func TooManyRequests(w http.ResponseWriter, r *http.Request) {
	v := map[string]any{"error": "too many requests"}
	response(w, r, http.StatusTooManyRequests, v)
}
