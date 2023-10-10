package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/romankravchuk/pastebin/config"
	v1 "github.com/romankravchuk/pastebin/internal/controller/http/v1"
	"github.com/romankravchuk/pastebin/pkg/httpserver"
	"github.com/romankravchuk/pastebin/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// HTTP Server
	handler := chi.NewMux()
	v1.NewRouter(handler)
	srv := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-srv.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := srv.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
