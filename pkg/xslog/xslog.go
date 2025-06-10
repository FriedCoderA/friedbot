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
	slog.Handler // 内嵌默认 Handler
	writer       io.Writer
}

func NewLogHandler(w io.Writer, level slog.Level) *LogHandler {
	// 使用 TextHandler 处理基础逻辑（包括 Enabled 检查）
	baseHandler := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
	})

	return &LogHandler{
		Handler: baseHandler,
		writer:  w,
	}
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	// 自定义日志格式
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
	UpdateLogLevel() // 初始化 handler
	return nil
}

func UpdateLogLevel() {
	level := getLogLevel()
	// 每次更新都创建新 Handler（确保级别生效）
	handler = NewLogHandler(multiWriter, level)
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
