package usecases

import (
	"context"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GetUserByIdInputDTO struct {
	ID string `json:"id"`
}

type GetUserByIdUseCase struct {
	userRepository entities.UserRepository
}

func NewGetUserByIdUseCase(userRepository entities.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{userRepository: userRepository}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context, input GetUserByIdInputDTO) (*entities.User, error) {
	log.Println("--> GET USUÁRIO PELO ID")
	user, err := uc.userRepository.GetByID(ctx, input.ID)
	if err != nil {
		log.Println("Erro ao buscar usuário pelo ID")
		return nil, err
	}
	log.Println("<- Dados do usuário lidos com sucesso")
	return user, nil
}
