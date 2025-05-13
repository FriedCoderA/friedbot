package main

import (
	"log"
	"log/slog"

	"friedbot/internal/config"
	"friedbot/internal/xslog"
)

func main() {
	if err := xslog.InitLog(); err != nil {
		log.Fatalf("Error initializing log: %v", err)
	}
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}
	slog.Info("start bot success")
}
