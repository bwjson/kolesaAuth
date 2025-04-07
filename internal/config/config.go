package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Env string
	Redis
	Twilio
	GRPC
}

type GRPC struct {
	Port int
}

type Redis struct {
	Address  string
	User     string
	Password string
}

type Twilio struct {
	AccountSid string
	AuthToken  string
	FromNumber string
}

func LoadConfig() *Config {
	godotenv.Load()

	var cfg Config

	cfg.Env = os.Getenv("ENV")

	cfg.Redis.Address = os.Getenv("REDIS_ADDRESS")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.User = os.Getenv("REDIS_USER")

	cfg.Twilio.AccountSid = os.Getenv("TWILIO_ACCOUNT_SID")
	cfg.Twilio.AuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	cfg.Twilio.FromNumber = os.Getenv("TWILIO_PHONE_NUMBER")

	cfg.GRPC.Port, _ = strconv.Atoi(os.Getenv("GRPC_PORT"))

	return &cfg
}
