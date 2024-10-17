package usecases

import "github.com/lucasBiazon/botany-back/internal/entities"

type GetByEmailInputDTO struct {
	Email string `json:"email"`
}

type GetByEmailUserOutputDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetByEmailUserUseCase struct {
	userRepo entities.UserRepository
}

func NewGetByEmailUserUseCase(repo entities.UserRepository) *GetByEmailUserUseCase {
	return &GetByEmailUserUseCase{userRepo: repo}
}

func (uc *GetByEmailUserUseCase) Execute(input GetByEmailInputDTO) (*GetByEmailUserOutputDTO, error) {
	user, err := uc.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	return &GetByEmailUserOutputDTO{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
