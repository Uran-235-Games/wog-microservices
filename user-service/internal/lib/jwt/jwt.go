package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"wog-server/internal/logger"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrSignedToken = errors.New("Ошибка подписи jwt токена")
)

type JWTSrv struct {
	secretKey string
	tokenDuration time.Duration
}

func NewJWTLib() *JWTSrv {
	secretKey := os.Getenv("JWT_TOKEN")
	if secretKey == "" {
		logger.Log.Warn("Secret JWT_TOKEN is not set in env. Setting JWT_TOKEN to default value")
		secretKey = "IvanRak"
	}
	tokenDurationStr := os.Getenv("JWT_DURATION")
	if tokenDurationStr == "" {
		logger.Log.Warn("Secret JWT_DURATION is not set in env. Setting JWT_DURATION to default value")
		tokenDurationStr = "1h"
	}
	tokenDuration, _ := time.ParseDuration(tokenDurationStr)
	return &JWTSrv{secretKey: secretKey, tokenDuration: tokenDuration}
}

func (s *JWTSrv) Generate(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["exp"] = time.Now().Add(s.tokenDuration).Unix()

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func (s *JWTSrv) Validate(tokenStr string) (userId string, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(string)
		if !ok {
			return "", fmt.Errorf("userID не найден или имеет неверный тип")
		}
		return id, nil
	}

	return "", fmt.Errorf("недействительный токен")
}
