package usecases_categorytask

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateCategoryTaskInputDTO struct {
	ID          string `json:"id"`
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

func (uc *UpdateCategoryTaskUseCase) Execute(ctx context.Context, input UpdateCategoryTaskInputDTO) error {
	log.Println("UpdateCategoryTaskUseCase - Execute")
	category, err := uc.CategoryTaskRepository.FindByID(ctx, input.UserID, input.ID)
	if err != nil {
		if err.Error() == "categoria de tarefa não encontrada" || err == nil {
			return errors.New("categoria de tarefa não encontrada")
		}
		return fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	if category.Name == input.Name && category.Description == input.Description && input.Name == "" && input.Description == "" {
		return errors.New("nenhum campo foi alterado")
	}
	if category.Name != input.Name && input.Name != "" {
		category.Name = input.Name
	}

	if category.Description != input.Description && input.Description != "" {
		category.Description = input.Description
	}

	if err := uc.CategoryTaskRepository.Update(ctx, category); err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	return nil
}
