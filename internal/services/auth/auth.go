package auth

import (
	"context"
	"github.com/bwjson/kolesa_auth/pkg/sms"
	"log/slog"
)

type Repository interface {
	Save(ctx context.Context, phoneNumber, verificationCode string) error
	Delete(ctx context.Context, phoneNumber string) error
	IsValidCode(ctx context.Context, phoneNumber, verificationCode string) (bool, error)
}

type SmsClient interface {
	SendSMS(ctx context.Context, phoneNumber string) error
}

type AuthService struct {
	log       *slog.Logger
	repo      Repository
	smsClient *sms.SmsClient
}

func NewAuthService(log *slog.Logger, repo Repository, smsClient *sms.SmsClient) *AuthService {
	return &AuthService{log: log, repo: repo, smsClient: smsClient}
}

type TokensResponse struct {
	accessToken  string
	refreshToken string
}

func (a *AuthService) SendVerificationCode(ctx context.Context, phoneNumber string) error {
	//Generate code
	//Save(),
	//SendSMS()
	return nil
}

func (a *AuthService) VerifyCode(ctx context.Context, phoneNumber, code string) (accessToken, refreshToken string, err error) {
	//IsValidCode()
	//JWT
	return "accessToken", "refreshToken", nil
}
