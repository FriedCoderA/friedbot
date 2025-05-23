package main

import (
	"log"
	"log/slog"

	"friedbot/internal/receiver"
	"friedbot/pkg/aigc"
	"friedbot/pkg/config"
	"friedbot/pkg/models"
	"friedbot/pkg/xslog"
)

func main() {
	if err := xslog.InitLog(); err != nil {
		log.Fatalf("Error initializing log: %v", err)
	}
	slog.Info("starting bot")

	slog.Info("reading configs")
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}

	slog.Info("initializing database")
	if err := models.InitModel(); err != nil {
		log.Fatalf("Error initializing model: %v", err)
	}

	slog.Info("initializing aigc client")
	if err := aigc.InitClient(); err != nil {
		log.Fatalf("Error initializing aigc client: %v", err)
	}

	slog.Info("initializing onebot service")
	if err := reply.NewService().Start(); err != nil {
		log.Fatalf("Error starting onebot service: %v", err)
	}
	slog.Info("start bot success")
}
