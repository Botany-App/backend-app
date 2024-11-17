package usecases_categorytask

import "github.com/lucasBiazon/botany-back/internal/entities"

type GetByNameCategoryTaskDTO struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type GetByNameCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewGetByNameCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *GetByNameCategoryTaskUseCase {
	return &GetByNameCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}

func (u *GetByNameCategoryTaskUseCase) Execute(dto *GetByNameCategoryTaskDTO) ([]entities.CategoryTask, error) {
	categoryTask, err := u.CategoryTaskRepository.GetByName(dto.UserID, dto.Name)
	if err != nil {
		return []entities.CategoryTask{}, err
	}

	return categoryTask, nil
}
