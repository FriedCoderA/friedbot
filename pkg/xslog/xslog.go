package xslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	logFilePath = "logs/friedbot.log"
)

var (
	logFile     *os.File
	multiWriter io.Writer
	handler     *LogHandler
)

type LogHandler struct {
	slog.Handler
	writer io.Writer
	level  slog.Level
}

func NewLogHandler(w io.Writer, level slog.Level) *LogHandler {
	return &LogHandler{
		writer: w,
		level:  level,
	}
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	logLine := fmt.Sprintf(
		"%s %-5s %s",
		r.Time.Format("2006-01-02 15:04:05"),
		"["+r.Level.String()+"]",
		r.Message,
	)
	r.Attrs(func(attr slog.Attr) bool {
		logLine += fmt.Sprintf(" %s=%v", attr.Key, attr.Value.String())
		return true
	})
	logLine += "\n"
	_, err := h.writer.Write([]byte(logLine))
	return err
}

func (h *LogHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level
}

func InitLog() error {
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if logFile == nil {
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		logFile = file
	}

	multiWriter = io.MultiWriter(os.Stdout, logFile)
	UpdateLogLevel() // 初始化时设置 handler
	return nil
}

func UpdateLogLevel() {
	level := getLogLevel()
	if handler == nil {
		handler = NewLogHandler(multiWriter, level)
	} else {
		handler.level = level
	}
	slog.SetDefault(slog.New(handler))
}

func getLogLevel() slog.Level {
	levelStr := viper.GetString("log.level")
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
