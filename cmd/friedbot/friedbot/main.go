package main

import (
	"log"
	"log/slog"

	"friedbot/internal/chat"
	"friedbot/internal/commands"
	"friedbot/internal/controllers"
	"friedbot/pkg/aigc"
	"friedbot/pkg/config"
	"friedbot/pkg/events"
	"friedbot/pkg/models"
	"friedbot/pkg/xslog"
)

func main() {
	if err := xslog.InitLog(); err != nil {
		log.Fatalf("initializing log failed: %v", err)
	}
	slog.Info("starting bot")

	slog.Info("reading configs")
	if err := config.InitConfig(); err != nil {
		log.Fatalf("initializing configuration failed: %v", err)
	}

	slog.Info("initializing database")
	if err := models.InitModel(); err != nil {
		log.Fatalf("initializing model failed: %v", err)
	}

	slog.Info("initializing aigc client")
	if err := aigc.InitClient(); err != nil {
		log.Fatalf("initializing aigc client failed: %v", err)
	}

	slog.Info("listening message events")
	if err := events.InitMessageEvents(); err != nil {
		log.Fatalf("listening message events failed: %v", err)
	}

	slog.Info("initializing command handlers")
	if err := commands.InitCommands(); err != nil {
		log.Fatalf("initializing command handlers failed: %v", err)
	}

	slog.Info("initializing chat bot")
	if err := chat.InitChatBot(); err != nil {
		log.Fatalf("initializing chat bot failed: %v", err)
	}

	slog.Info("starting server")
	if err := controllers.NewMainController().Start(); err != nil {
		log.Fatalf("starting server failed: %v", err)
	}

	slog.Info("start bot success")
}
