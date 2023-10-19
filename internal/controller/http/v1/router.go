// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/romankravchuk/pastebin/config"
	"github.com/romankravchuk/pastebin/internal/controller/http/v1/auth"
	"github.com/romankravchuk/pastebin/internal/controller/http/v1/paste"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/internal/usecase/repo"
	"github.com/romankravchuk/pastebin/pkg/log"
	"github.com/romankravchuk/pastebin/pkg/minio"
	"github.com/romankravchuk/pastebin/pkg/postgres"
	swagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/romankravchuk/pastebin/docs" //
)

// NewRouter returns a new router for api v1.
//
// Swagger spec:
//
//	@title						Pastebin API
//	@description				Implementation pastebin API
//	@version					1.0
//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securitydefinitions.apiKey	Bearer
//	@in							header
//	@name						Authorization
func NewRouter(mux chi.Router, cfg *config.Config, l *log.Logger) error {
	pg, err := postgres.New(cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	m, err := minio.New(cfg.Minio.DSN, cfg.Minio.AccessKey, cfg.Minio.SecretKey, 10*time.Second)
	if err != nil {
		return err
	}

	var (
		prepo = repo.NewPastesRepo(pg, m)
		auc   = usecase.NewAuth(nil)
		puc   = usecase.NewPastes(prepo, nil)
	)

	mux.Get("/swagger/*", swagger.Handler())

	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { render.Status(r, http.StatusOK) })

	mux.Method(http.MethodGet, "/metrics", promhttp.Handler())

	auth.New(mux, auc, l)

	paste.New(mux, puc, l)

	return nil
}
