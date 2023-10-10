// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/romankravchuk/pastebin/docs"
	swagger "github.com/swaggo/http-swagger/v2"
)

// NewRouter
//
// Swagger spec:
//
//	@title			Pastebin API
//	@description	Implementation pastebin API
//	@version		1.0
//	@host			localhost:8080
//	@BasePath		/v1
func NewRouter(mux *chi.Mux) {
	mux.Get("/swagger/*", swagger.Handler())

	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { render.Status(r, http.StatusOK) })

	mux.Method(http.MethodGet, "/metrics", promhttp.Handler())
}
