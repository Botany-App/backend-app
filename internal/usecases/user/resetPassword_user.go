package usecases

import (
	"context"
	"errors"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
)

type ResetPasswordUserInputDTO struct {
	newPassword string `json:"newPassword"`
	token       string `json:"token"`
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
	userID, err := services.ExtractUserIDFromToken(input.token, uc.JWTService)
	if err != nil {
		return err
	}

	log.Println("--> Pegando dados de usuário")
	user, err := uc.UserRepository.GetByID(ctx, userID)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	user.Password = input.newPassword
	log.Println("--> Atualizando senha de usuário")
	err = uc.UserRepository.UpdatePassword(ctx, user.ID.String(), user.Password)
	if err != nil {
		return errors.New("erro ao atualizar senha de usuário")
	}

	return nil
}
