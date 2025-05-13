package main

import (
	"log"
	"log/slog"

	"friedbot/internal/config"
	"friedbot/internal/xslog"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}
	if err := xslog.InitLog(); err != nil {
		log.Fatalf("Error initializing log: %v", err)
	}
	slog.Info("Starting bot...")
	slog.Info("Starting bot...")
}
