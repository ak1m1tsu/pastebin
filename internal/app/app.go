package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/romankravchuk/pastebin/config"
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
	v1 "github.com/romankravchuk/pastebin/internal/controller/http/v1"
	"github.com/romankravchuk/pastebin/pkg/httpserver"
	"github.com/romankravchuk/pastebin/pkg/log"
)

func Run(cfg *config.Config) {
	var (
		err error
		l   = log.New(os.Stdout, log.Stol(cfg.Log.Level))
	)

	// HTTP Server
	handler := chi.NewMux()
	handler.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.NotFound(w, r)
	})
	handler.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.MethodNotAllowed(w, r)
	})
	handler.Route("/api/v1", func(r chi.Router) {
		err = v1.NewRouter(r, cfg, l)
	})

	if err != nil {
		l.Error("initialize the router", err, nil)

		return
	}

	srv := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	l.Info("start the service", log.FF{
		{Key: "name", Value: cfg.App.Name},
		{Key: "version", Value: cfg.App.Version},
		{Key: "port", Value: cfg.HTTP.Port},
	})

	select {
	case s := <-interrupt:
		l.Info("catch shutdown signal", log.FF{{
			Key:   "signal",
			Value: s.String(),
		}})
	case err = <-srv.Notify():
		l.Error("catch listen and serve notification", err, nil)
	}

	// Shutdown
	err = srv.Shutdown()
	if err != nil {
		l.Error("shutdown the service", err, nil)
	}
}
