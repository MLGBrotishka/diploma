package jwt

import (
	"fmt"
	"time"

	"auth/internal/entity"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
}

func New(secret string) *JWTService {
	return &JWTService{
		secret: secret,
	}
}

// NewToken создает новый JWT токен для пользователя.
func (s *JWTService) NewToken(user entity.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) ParseToken(token string) (entity.User, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to parse token: %w", err)
	}

	user := entity.User{
		ID:    claims["uid"].(int64),
		Login: claims["login"].(string),
	}

	return user, nil
}
