package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/romankravchuk/pastebin/config"
	v1 "github.com/romankravchuk/pastebin/internal/controller/http/v1"
	"github.com/romankravchuk/pastebin/pkg/httpserver"
	"github.com/romankravchuk/pastebin/pkg/log"
)

func Run(cfg *config.Config) {
	l := log.New(os.Stdout, log.Stol(cfg.Log.Level))

	// HTTP Server
	handler := chi.NewMux()
	v1.NewRouter(handler)
	srv := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	l.Info("app - Run - starting service", log.FF{
		{Key: "name", Value: cfg.App.Name},
		{Key: "version", Value: cfg.App.Version},
		{Key: "port", Value: cfg.HTTP.Port},
	})

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal", log.FF{{
			Key:   "signal",
			Value: s.String(),
		}})
	case err := <-srv.Notify():
		l.Error("app - Run - httpServer.Notify", err, nil)
	}

	// Shutdown
	err := srv.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown", err, nil)
	}
}
