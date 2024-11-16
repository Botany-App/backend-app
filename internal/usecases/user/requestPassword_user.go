package usecases

import (
	"context"
	"errors"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
)

type RequestPasswordResetUserUseCase struct {
	UserRepository entities.UserRepository
	JWTService     services.JWTService
}

type RequestPasswordResetUserInputDTO struct {
	Email string `json:"email"`
}

func NewRequestPasswordResetUseCase(userRepo entities.UserRepository, jwtService services.JWTService) *RequestPasswordResetUserUseCase {
	return &RequestPasswordResetUserUseCase{
		UserRepository: userRepo,
		JWTService:     jwtService,
	}
}

func (uc *RequestPasswordResetUserUseCase) Execute(ctx context.Context, email string) error {
	user, err := uc.UserRepository.GetByEmail(ctx, email)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	resetToken, err := uc.JWTService.GenerateToken(user.ID.String())
	if err != nil {
		return errors.New("erro ao gerar token de redefinição de senha")
	}

	err = services.NewEmailService().SendEmailResetPassword(user.Email, resetToken)
	if err != nil {
		return errors.New("erro ao enviar email de redefinição de senha")
	}
	return nil
}
