package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/adal4ik/crypto-service/internal/config"
	"github.com/adal4ik/crypto-service/internal/handler"
	"github.com/adal4ik/crypto-service/internal/repository"
	"github.com/adal4ik/crypto-service/internal/service"
	"github.com/adal4ik/crypto-service/pkg/loadenv"
	"github.com/adal4ik/crypto-service/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	loadenv.LoadEnvFile(".env")
	logger := logger.New(os.Getenv("APP_ENV"))
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "Error syncing logger: %v\n", err)
		}
	}()
	logger.Info("Starting Subtracker application", zap.String("environment", os.Getenv("APP_ENV")))
	// Load configuration
	cfg := config.LoadConfig()
	logger.Info("Configuration loaded", zap.String("app_port", cfg.App.AppPort), zap.String("log_level", cfg.App.LogLevel),
		zap.String("db_host", cfg.Postgres.DBHost), zap.String("db_port", cfg.Postgres.DBPort),
		zap.String("db_name", cfg.Postgres.DBName), zap.String("db_user", cfg.Postgres.DBUser),
		zap.String("db_password", cfg.Postgres.DBPassword), zap.String("postgres_dsn", cfg.Postgres.PostgresDSN))

	db, err := repository.ConnectDB(ctx, cfg.Postgres, logger)
	if err != nil {
		logger.Fatal("Failed to connect to the database", zap.Error(err))
	}
	defer db.Close()
	logger.Info("Connected to the database successfully", zap.String("dsn", cfg.Postgres.PostgresDSN))
	repo := repository.NewRepository(db, logger)
	service := service.NewService(repo, logger)
	handlers := handler.NewHandlers(service, logger)
	logger.Info("All components initialized successfully")
	mux := handler.Router(*handlers)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		log.Println("Server is running on port: http://localhost" + httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe error", zap.Error(err))
		}
	}()
	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("HTTP server shutdown error", zap.Error(err))
	}

	logger.Info("Server stopped gracefully")

}
