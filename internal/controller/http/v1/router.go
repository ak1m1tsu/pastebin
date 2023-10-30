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
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
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
	postgreClient, err := postgres.New(cfg.Postgres.DSN)
	if err != nil {
		return err
	}

	minioClient, err := minio.New(cfg.Minio.DSN, cfg.Minio.AccessKey, cfg.Minio.SecretKey)
	if err != nil {
		return err
	}

	redisClient, err := redis.New(cfg.Redis.DSN)
	if err != nil {
		return err
	}

	var (
		pastesCache   = cache.NewPastesCache(redisClient)
		pastesBlob    = blob.NewPastesBlobStorage(minioClient)
		pastesRepo    = repo.NewPastesRepositry(postgreClient)
		usersRepo     = repo.NewUsersRepositry(postgreClient)
		oauthapi      = webapi.NewGithubAPI(cfg.OAuth.ClientID, cfg.OAuth.ClientSecret)
		authUsecase   = usecase.NewAuth(usersRepo, oauthapi)
		pastesUsecase = usecase.NewPastes(pastesRepo, pastesBlob, pastesCache)
	)

	mux.Use(middleware.RedirectSlashes)
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

	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { response.OK(w, r, render.M{"status": "alive"}) })

	mux.Method(http.MethodGet, "/metrics", promhttp.Handler())

	auth.MountRoutes(mux, authUsecase, l)

	paste.MountRoutes(mux, pastesUsecase, l)

	return nil
}
