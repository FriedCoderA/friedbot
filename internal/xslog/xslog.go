package xslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	logFilePath = "log/friedbot.log"
)

type LogHandler struct {
	slog.Handler
	writer io.Writer
	opts   slog.HandlerOptions
}

func NewLogHandler(w io.Writer, opts *slog.HandlerOptions) *LogHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &LogHandler{
		writer: w,
		opts:   *opts,
	}
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	logLine := fmt.Sprintf(
		"%s %-5s %s",
		r.Time.Format("2006-01-02 15:04:05"),
		"["+r.Level.String()+"]",
		r.Message,
	)
	_, err := h.writer.Write([]byte(logLine))
	return err
}

func (h *LogHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.opts.Level.Level()
}

func InitLog() error {
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	handler := NewLogHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
	return nil
}
