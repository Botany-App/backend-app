package usecases

import (
	"context"
	"log"
	"os"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
)

type LoginUserUseCase struct {
	UserRepository entities.UserRepository
}

type LoginUserInputDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewLoginUserUseCase(userRepo entities.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		UserRepository: userRepo,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	log.Println("--lOGANDO---")
	id, err := uc.UserRepository.Login(ctx, email, password)
	if err != nil {
		if id == "not found" {
			return "user not found", err
		}
		if id == "invalid password" {
			return "invalid password", err
		}
		return "", err
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenJWT, err := services.NewJWTService(secretKey).GenerateToken(id)
	if err != nil {
		return "", err
	}

	return tokenJWT, nil
}
