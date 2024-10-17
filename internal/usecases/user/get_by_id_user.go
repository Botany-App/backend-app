package usecases

import "github.com/lucasBiazon/botany-back/internal/entities"

type GetByIdInputDTO struct {
	ID string `json:"id"`
}

type GetByIdUserOutputDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GetByIdUserUseCase struct {
	userRepo entities.UserRepository
}

func NewGetByIdUserUseCase(repo entities.UserRepository) *GetByIdUserUseCase {
	return &GetByIdUserUseCase{userRepo: repo}
}

func (uc *GetByIdUserUseCase) Execute(input GetByIdInputDTO) (*GetByIdUserOutputDTO, error) {
	user, err := uc.userRepo.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetByIdUserOutputDTO{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}
