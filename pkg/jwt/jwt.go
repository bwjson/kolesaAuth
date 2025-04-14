package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strings"
	"time"
)

type JWTClient struct {
	JWTSecret string
}

func NewJWTClient(JWTSecret string) *JWTClient {
	return &JWTClient{JWTSecret: JWTSecret}
}

type tokenClaims struct {
	jwt.StandardClaims
	PhoneNumber string `json:"phone_number"`
}

func (c *JWTClient) GenerateAccessToken(phoneNumber string) (string, error) {
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"phone_number": phoneNumber,
			"exp":          time.Now().Add(time.Minute * 15).Unix(),
		},
	)

	log.Println(c.JWTSecret)
	tokenString, err := tokenWithClaims.SignedString([]byte(c.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *JWTClient) GenerateRefreshToken(phoneNumber string) (string, error) {
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"phone_number": phoneNumber,
			"exp":          time.Now().Add(time.Hour * 24 * 31).Unix(),
		},
	)

	tokenString, err := tokenWithClaims.SignedString([]byte(c.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *JWTClient) ParseToken(tokenString string) (string, error) {
	// This function check signature and return phone number from the payload part of token
	tokenString = strings.TrimSpace(tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		log.Println(c.JWTSecret)
		return []byte(c.JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("Invalid token claims")
	}

	return claims.PhoneNumber, nil
}
