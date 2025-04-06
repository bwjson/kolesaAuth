package app

import (
	"context"
	grpcapp "github.com/bwjson/kolesa_auth/internal/app/grpc"
	"github.com/bwjson/kolesa_auth/internal/redis"
	"github.com/bwjson/kolesa_auth/internal/services/auth"
	"github.com/bwjson/kolesa_auth/pkg/sms"
	"github.com/twilio/twilio-go"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	addr, user, pass, accountSID, authToken, fromNumber string,
) *App {
	ctx := context.Context(context.Background())

	cache := redis.NewRedisClient(ctx, addr, user, pass)

	client := twilio.NewRestClient()

	smsClient := sms.NewSmsClient(client, accountSID, authToken, fromNumber)

	repo := redis.NewRepository(cache)

	authService := auth.NewAuthService(log, repo, smsClient)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
