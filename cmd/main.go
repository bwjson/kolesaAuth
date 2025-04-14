package main

import (
	"github.com/bwjson/kolesa_auth/internal/app"
	"github.com/bwjson/kolesa_auth/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	log := setupLogger(cfg.Env)

	log.Info("Starting...")

	application := app.New(log,
		cfg.GRPC.Port, cfg.Redis.Address, cfg.Redis.User, cfg.Redis.Password,
		cfg.Twilio.AccountSid, cfg.Twilio.AuthToken, cfg.Twilio.FromNumber, cfg.JWT.JWTSecret)

	go func() {
		application.GRPCServer.MustRun()
	}()

	log.Info("gRPC server started", "port", cfg.GRPC.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
