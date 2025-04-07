package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/bwjson/kolesa_auth/internal/lib/random/codeutil"
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

func (a *AuthService) SendVerificationCode(ctx context.Context, phoneNumber string) error {
	code, err := codeutil.GenerateFourDigitsCode()
	if err != nil {
		return err
	}

	err = a.repo.Save(ctx, phoneNumber, code)
	if err != nil {
		return err
	}

	err = a.smsClient.SendSMS(fmt.Sprintf("Your verification code is %s", code), phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) VerifyCode(ctx context.Context, phoneNumber, code string) (accessToken, refreshToken string, err error) {
	a.log.Info("Starting the service method VerifyCode")

	isValid, err := a.repo.IsValidCode(ctx, phoneNumber, code)

	a.log.Info("IsValidCode method used")

	if err != nil {
		return "", "", err
	}
	if !isValid {
		return "", "", errors.New("invalid code")
	}

	//JWT

	a.repo.Delete(ctx, phoneNumber)

	return "accessToken", "refreshToken", nil
}
