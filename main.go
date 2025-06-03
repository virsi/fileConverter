package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/virsi/fileConverter/internal/config"
	"github.com/virsi/fileConverter/internal/storage/sqlite"
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

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("Storage is not init", slog.String("error", err.Error())) // TODO: create err func in internal/lib/logger/slog -- slog.Attr{key: error, value: err.Error()} slog.StringValue(err.Error())
		os.Exit(1)
	}
	_ = storage

	id, err := storage.SaveFile("example.jpg", "jpg", "png", "new", "/path/to/file.jpg")
	if err != nil {
		log.Error("Failed to save file", slog.String("error", err.Error()))
	}
	log.Info("File saved successfully", slog.Int64("id", id))

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
