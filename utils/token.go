package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/youtube/config"
)

func Generate_JWT_Token(userId, email string) (string, error) {
	cfg, err := config.SetUpConfig()
	if err != nil {
		return "", err
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userId
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(cfg.JWT_TOKEN))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func Verify_JWT_Token(tokenString string) (*jwt.Token, error) {
	cfg, err := config.SetUpConfig()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
