package config

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type loggerCtxKey struct{}

func NewLogger(appEnv string) *slog.Logger {
	_ = os.MkdirAll("logs", 0o755)

	fileRotator := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    20,   // Ukuran maksimal per file log dalam MB
		MaxBackups: 5,    // Maksimal jumlah file backup lama yang disimpan
		MaxAge:     14,   // Retensi penyimpanan file log lama dalam hari
		Compress:   true, // Kompres file log lama ke .gz
	}

	multiWriter := io.MultiWriter(os.Stdout, fileRotator)

	level := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(appEnv)) {
	case "development", "test":
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(loggerCtxKey{}).(*slog.Logger); ok && l != nil {
			return l
		}
	}
	return slog.Default()
}