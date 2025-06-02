package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virsi/fileConverter/internal/config"
)

const (
	envLocal = "local"
	envDevelopment = "development"
	envProduction = "production"
	envTest = "test"
)

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting file converter service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO init storage
	// TODO init router
	// TODO start server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDevelopment:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
