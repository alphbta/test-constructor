package utils

import (
	"test-constructor/config"
	"test-constructor/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userID uint, email, name, surname string, role int) (string, error) {
	cfg := config.Load()

	claims := &models.JWTClaims{
		UserID:  userID,
		Email:   email,
		Name:    name,
		Surname: surname,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.JWTTTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test-constructor",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	cfg := config.Load()

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
