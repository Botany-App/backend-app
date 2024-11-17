package usecases_categorytask

import (
	"errors"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateCategoryTaskDTO struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewUpdateCategoryTaskUseCase(repository entities.CategoryTaskRepository) *UpdateCategoryTaskUseCase {
	return &UpdateCategoryTaskUseCase{CategoryTaskRepository: repository}
}

func (uc UpdateCategoryTaskUseCase) Execute(dto UpdateCategoryTaskDTO) error {
	categories, err := uc.CategoryTaskRepository.GetByID(dto.UserID, dto.Name)
	if err != nil {
		return err
	}

	category := categories[0]
	if category.Name != dto.Name && dto.Name != "" {
		category.Name = dto.Name
	}

	if category.Description != dto.Description && dto.Description != "" {
		category.Description = dto.Description
	}

	if category.Name == dto.Name && category.Description == dto.Description && dto.Name == "" && dto.Description == "" {
		return errors.New("nenhum campo foi alterado")
	}

	err = uc.CategoryTaskRepository.Update(dto.UserID, &category)
	if err != nil {
		return err
	}

	return nil
}
