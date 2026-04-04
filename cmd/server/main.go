package main

import (
	"context"
	"fatearcan/internal/config"
	"fatearcan/internal/handlers"
	"fatearcan/internal/services"
	"fatearcan/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// 1. Загрузка конфига
	cfg := config.MustLoad()

	// 2. Инициализация логгера
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	// 3. Инициализация хранилища (карты и расклады)
	store, err := storage.NewStorage(cfg.Storage.DeckPath, cfg.Storage.SpreadPath)
	if err != nil {
		log.Error("failed to init storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("storage initialized")

	// 4. Выбор AI провайдера
	var aiProvider services.AIPredictor

	switch cfg.AI.Mode {
	case "cloud":
		aiProvider = services.NewCloudPredictor(
			log,
			cfg.AI.Cloud.URL,
			cfg.AI.Cloud.Key,
			cfg.AI.Cloud.FolderID,
			cfg.AI.Cloud.Model,
		)
		log.Info("using CLOUD AI provider (Yandex)")
	case "local":
		aiProvider = services.NewLocalPredictor(
			log,
			cfg.AI.Local.URL,
			cfg.AI.Local.Model,
		)
		log.Info("using LOCAL AI provider (Ollama)")
	default:
		log.Error("unknown AI mode", slog.String("mode", cfg.AI.Mode))
		os.Exit(1)
	}

	// 5. Инициализация сервиса
	tarotService := services.NewTarotService(store, aiProvider)

	// 6. Инициализация хендлеров
	tarotHandler := handlers.NewTarotHandler(tarotService, log)

	// 7. Роутинг
	router := http.NewServeMux()
	router.HandleFunc("GET /api/spreads", tarotHandler.HandleSpreads)
	router.HandleFunc("POST /api/reading", tarotHandler.HandleReading)
	router.HandleFunc("POST /api/chat", tarotHandler.HandleChat)
	router.Handle("/", http.FileServer(http.Dir("./web")))

	// 8. Запуск сервера
	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", slog.String("error", err.Error()))
		}
	}()

	log.Info("server started", slog.String("address", cfg.HttpServer.Address))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", slog.String("error", err.Error()))
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		// Default fallback
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
