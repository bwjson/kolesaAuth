package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwjson/kolesa_auth/internal/lib/random/codeutil"
	"github.com/bwjson/kolesa_auth/pkg/bot"
	"github.com/bwjson/kolesa_auth/pkg/jwt"
	"github.com/bwjson/kolesa_auth/pkg/sms"
	"log/slog"
	"time"
)

type Repository interface {
	Save(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
	IsValidCode(ctx context.Context, phoneNumber, verificationCode string) (bool, error)
	GetAll(ctx context.Context) (map[string]string, error)
}

type SmsClient interface {
	SendSMS(ctx context.Context, phoneNumber string) error
}

type AuthService struct {
	log       *slog.Logger
	repo      Repository
	smsClient *sms.SmsClient
	jwtClient *jwt.JWTClient
	botClient *bot.BotClient
}

func NewAuthService(log *slog.Logger, repo Repository, smsClient *sms.SmsClient, jwtClient *jwt.JWTClient, botClient *bot.BotClient) *AuthService {
	return &AuthService{log: log, repo: repo, smsClient: smsClient, jwtClient: jwtClient, botClient: botClient}
}

func (a *AuthService) SendVerificationCode(ctx context.Context, phoneNumber string) error {
	code, err := codeutil.GenerateFourDigitsCode()
	if err != nil {
		return err
	}

	err = a.repo.Save(ctx, "ACCESS"+phoneNumber, code, time.Minute*3)
	if err != nil {
		return err
	}

	err = a.botClient.SendMessageToChannel(fmt.Sprintf("Phone: %s, Your verification code is: %v", phoneNumber, code))
	if err != nil {
		return err
	}

	// Production SMS
	//err = a.smsClient.SendSMS(fmt.Sprintf("Your verification code is %s", code), phoneNumber)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (a *AuthService) VerifyCode(ctx context.Context, phoneNumber, code string) (accessToken, refreshToken string, err error) {
	isValid, err := a.repo.IsValidCode(ctx, "ACCESS"+phoneNumber, code)

	if err != nil {
		return "", "", err
	}
	if !isValid {
		return "", "", errors.New("invalid code")
	}

	// Generating JWT tokens
	accessToken, err = a.jwtClient.GenerateAccessToken(phoneNumber)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = a.jwtClient.GenerateRefreshToken(phoneNumber)
	if err != nil {
		return "", "", err
	}

	// ttl is one week
	err = a.repo.Save(ctx, "REFRESH"+phoneNumber, refreshToken, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// Parse payload from token
	phoneNumber, err := a.jwtClient.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Check if this token in cache
	_, err = a.repo.Get(ctx, "REFRESH"+phoneNumber)
	if err != nil {
		return "", err
	}

	accessToken, err := a.jwtClient.GenerateAccessToken(phoneNumber)
	if err != nil {
		return "", err
	}

	a.repo.Save(ctx, "ACCESS"+phoneNumber, accessToken, time.Minute*3)

	// Refresh token remains the same
	return accessToken, nil
}
