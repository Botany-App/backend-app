package usecases

import (
	"context"
	"log"
	"os"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
)

type LoginUserInputDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserUseCase struct {
	UserRepository entities.UserRepository
}

func NewLoginUserUseCase(userRepo entities.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		UserRepository: userRepo,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, input LoginUserInputDTO) (string, error) {
	log.Println("LoginUserUseCase - Execute")
	ID, err := uc.UserRepository.Login(ctx, input.Email, input.Password)
	if err != nil {
		if ID == "not found" {
			log.Println("User not found")
		}
		if ID == "invalid password" {
			log.Println("Invalid password")
		}
		return "", err
	}

	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenJWT, err := services.NewJWTService(secretKey).GenerateToken(ID)
	if err != nil {
		return "", err
	}
	return tokenJWT, nil
}
