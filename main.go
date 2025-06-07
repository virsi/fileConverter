package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/virsi/fileConverter/internal/config"
	"github.com/virsi/fileConverter/internal/http-server/middleware/logger"
	"github.com/virsi/fileConverter/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDevelopment = "development"
	envProduction = "production"
	envTest = "test"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()
	fmt.Print(cfg)

	// Setup logger based on environment
	// TODO: setup a pretty logger
	log := setupLogger(cfg.Env)
	log.Info("starting file converter service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Initialize storage
	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("Storage is not init", slog.String("error", err.Error())) // TODO: create err func in internal/lib/logger/slog -- slog.Attr{key: error, value: err.Error()} slog.StringValue(err.Error())
		os.Exit(1)
	}
	_ = storage

	// Initialize router
	// TODO: add cors middleware
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat) // works only with chi
	router.Use(mwlogger.New(log))

	// GetFileByID test
	// file, err := storage.GetFileByID(1)
	// if err != nil {
	// 	log.Error("Failed to get file by ID", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }
	// log.Info("File retrieved", slog.String("file", file["original_filename"]))

	// UpdateFileStatus test
	// err = storage.UpdateFileStatus(1, "processing")
	// if err != nil {
	// 	log.Error("Failed to update file status", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }
	// log.Info("File status updated", slog.Int64("id", 1), slog.String("status", "processing"))

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
