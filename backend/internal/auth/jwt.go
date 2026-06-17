package auth

import (
	"test-constructor/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userID uint, email, name, surname, role string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	config *config.Config
}

func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{config: cfg}
}

func (s *jwtService) GenerateToken(userID uint, email, name, surname, role string) (string, error) {
	claims := &JWTClaims{
		UserID:  userID,
		Email:   email,
		Name:    name,
		Surname: surname,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.config.JWTTTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test-constructor",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
