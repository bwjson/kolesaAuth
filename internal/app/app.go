package app

import (
	"context"
	grpcapp "github.com/bwjson/kolesa_auth/internal/app/grpc"
	"github.com/bwjson/kolesa_auth/internal/redis"
	"github.com/bwjson/kolesa_auth/internal/services/auth"
	bot "github.com/bwjson/kolesa_auth/pkg/bot"
	"github.com/bwjson/kolesa_auth/pkg/jwt"
	"github.com/bwjson/kolesa_auth/pkg/sms"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/twilio/twilio-go"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	addr, user, pass, accountSID, authToken, fromNumber, jwtSecret, botToken, channelID string,
) *App {
	ctx := context.Context(context.Background())

	cache := redis.NewRedisClient(ctx, addr, user, pass)

	jwtClient := jwt.NewJWTClient(jwtSecret)

	client := twilio.NewRestClient()

	smsClient := sms.NewSmsClient(client, accountSID, authToken, fromNumber)

	tgBot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	botClient := bot.NewBotClient(tgBot, botToken, channelID)

	repo := redis.NewRepository(cache)

	authService := auth.NewAuthService(log, repo, smsClient, jwtClient, botClient)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
