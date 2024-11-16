package usecases

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
)

type RequestPasswordResetUserInputDTO struct {
	Email string `json:"email"`
}

type RequestPasswordResetUserUseCase struct {
	UserRepository entities.UserRepository
	JWTService     services.JWTService
}

func NewRequestPasswordResetUseCase(userRepo entities.UserRepository, jwtService services.JWTService) *RequestPasswordResetUserUseCase {
	return &RequestPasswordResetUserUseCase{
		UserRepository: userRepo,
		JWTService:     jwtService,
	}
}

func (uc *RequestPasswordResetUserUseCase) Execute(ctx context.Context, input RequestPasswordResetUserInputDTO) error {
	log.Println("-> Request Password Reset User")
	user, err := uc.UserRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	log.Println("-> Generate Token")
	resetToken, err := uc.JWTService.GenerateToken(user.ID)
	if err != nil {
		return errors.New("erro ao gerar token de redefinição de senha")
	}

	log.Println("-> Send Email")
	err = services.NewEmailService().SendEmailResetPassword(user.Email, resetToken)
	if err != nil {
		return errors.New("erro ao enviar email de redefinição de senha")
	}
	return nil
}
