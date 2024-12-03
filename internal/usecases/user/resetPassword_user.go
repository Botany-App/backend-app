package usecases

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordUserInputDTO struct {
	NewPassword string `json:"newPassword"`
	Token       string `json:"token"`
}

type ResetPasswordUserUseCase struct {
	UserRepository entities.UserRepository
	JWTService     services.JWTService
}

func NewResetPasswordUserUseCase(userRepository entities.UserRepository, jwtService services.JWTService) *ResetPasswordUserUseCase {
	return &ResetPasswordUserUseCase{
		UserRepository: userRepository,
		JWTService:     jwtService,
	}
}

func (uc *ResetPasswordUserUseCase) Execute(ctx context.Context, input ResetPasswordUserInputDTO) error {
	log.Println("ResetPasswordUserUseCase - Execute")
	if uc.UserRepository.IsTokenRevokedPassword(ctx, input.Token) {
		return errors.New("token inválido ou já utilizado")
	}

	userID, err := services.ExtractUserIDFromToken(input.Token, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return err
	}

	user, err := uc.UserRepository.FindByID(ctx, userID)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao gerar hash de senha")
	}

	user.Password = string(passwordHash)

	err = uc.UserRepository.UpdatePassword(ctx, user.Id, user.Password)
	if err != nil {
		return errors.New("erro ao atualizar senha de usuário")
	}

	err = uc.UserRepository.StoreRevokedTokenPassword(ctx, input.Token)
	if err != nil {
		return err
	}

	return nil
}
