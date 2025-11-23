package models

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Role    int    `json:"role"`
	jwt.RegisteredClaims
}
