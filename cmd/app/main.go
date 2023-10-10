package main

import (
	"log"

	"github.com/romankravchuk/pastebin/config"
	"github.com/romankravchuk/pastebin/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
