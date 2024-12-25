package helpers

import (
	"event-management/structs"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Claims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func SignPayload(payload structs.User) (string, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		ID:   payload.ID,
		Role: payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(viper.GetString("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
