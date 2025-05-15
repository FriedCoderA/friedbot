package main

import (
	"log"
	"log/slog"

	"friedbot/pkg/config"
	"friedbot/pkg/models"
	"friedbot/pkg/xslog"
)

func main() {
	if err := xslog.InitLog(); err != nil {
		log.Fatalf("Error initializing log: %v", err)
	}
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}
	if err := models.InitModel(); err != nil {
		log.Fatalf("Error initializing model: %v", err)
	}
	slog.Info("start bot success")
}
