package usecases_categorytask

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateCategoryTaskDTO struct {
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

func (uc *UpdateCategoryTaskUseCase) Execute(ctx context.Context, input UpdateCategoryTaskDTO) error {
	log.Println("--> Get category by ID")
	category, err := uc.CategoryTaskRepository.GetByID(ctx, input.UserID, input.ID)
	if err != nil {
		// Ajuste para verificar o erro de forma mais genérica
		if err.Error() == "categoria de tarefa não encontrada" || err == nil {
			return errors.New("categoria de tarefa não encontrada")
		}
		return fmt.Errorf("erro ao buscar categoria: %w", err)
	}

	log.Println("-> Verificando campos alterados")
	if category.Name == input.Name && category.Description == input.Description && input.Name == "" && input.Description == "" {
		return errors.New("nenhum campo foi alterado")
	}
	if category.Name != input.Name && input.Name != "" {
		category.Name = input.Name
	}

	if category.Description != input.Description && input.Description != "" {
		category.Description = input.Description
	}

	log.Println("--> Atualizando categoria")
	if err := uc.CategoryTaskRepository.Update(ctx, category); err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %w", err)
	}

	return nil
}
