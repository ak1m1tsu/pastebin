// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/romankravchuk/pastebin/config"
	"github.com/romankravchuk/pastebin/internal/controller/http/middleware/logger"
	"github.com/romankravchuk/pastebin/internal/controller/http/v1/auth"
	"github.com/romankravchuk/pastebin/internal/controller/http/v1/paste"
	"github.com/romankravchuk/pastebin/internal/usecase"
	"github.com/romankravchuk/pastebin/internal/usecase/blob"
	"github.com/romankravchuk/pastebin/internal/usecase/cache"
	"github.com/romankravchuk/pastebin/internal/usecase/repo"
	"github.com/romankravchuk/pastebin/internal/usecase/webapi"
	"github.com/romankravchuk/pastebin/pkg/log"
	"github.com/romankravchuk/pastebin/pkg/minio"
	"github.com/romankravchuk/pastebin/pkg/postgres"
	"github.com/romankravchuk/pastebin/pkg/redis"
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

	m, err := minio.New(cfg.Minio.DSN, cfg.Minio.AccessKey, cfg.Minio.SecretKey)
	if err != nil {
		return err
	}

	rd, err := redis.New(cfg.Redis.DSN)
	if err != nil {
		return err
	}

	var (
		pcache   = cache.NewPastesCache(rd)
		pblob    = blob.New(m)
		prepo    = repo.NewPastesRepo(pg)
		urepo    = repo.NewUsersRepo(pg)
		oauthapi = webapi.New(cfg.OAuth.ClientID, cfg.OAuth.ClientSecret)
		auc      = usecase.NewAuth(urepo, oauthapi)
		puc      = usecase.NewPastes(prepo, pblob, pcache)
	)

	mux.Use(middleware.RealIP)
	mux.Use(logger.New(l))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Link", "Location"},
		MaxAge:           300,
	}))

	mux.Get("/swagger/*", swagger.Handler())

	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { render.Status(r, http.StatusOK) })

	mux.Method(http.MethodGet, "/metrics", promhttp.Handler())

	auth.New(mux, auc, l)

	paste.New(mux, puc, l)

	return nil
}
