package services

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userID interface{}) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type JWTServiceImpl struct {
	secretKey string
}

func NewJWTService(secretKey string) *JWTServiceImpl {
	return &JWTServiceImpl{
		secretKey: secretKey,
	}
}

func (j *JWTServiceImpl) GenerateToken(userID interface{}) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verificar método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secretKey), nil
	})
}

func ExtractUserIDFromToken(tokenString string, jwtService JWTService) (string, error) {
	if tokenString == "" {
		return "", errors.New("token não fornecido")
	}

	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.Split(tokenString, " ")[1]
	}

	token, err := jwtService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return "", errors.New("token inválido ou expirado")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("falha ao obter claims do token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("userID não encontrado no token")
	}

	return userID, nil
}
