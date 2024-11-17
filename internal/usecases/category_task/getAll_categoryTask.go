package usecases_categorytask

import "github.com/lucasBiazon/botany-back/internal/entities"

type GetAllCategoryTaskDTO struct {
	UserID string `json:"user_id"`
}

type GetAllCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

func NewGetAllCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *GetAllCategoryTaskUseCase {
	return &GetAllCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}

func (uc *GetAllCategoryTaskUseCase) Execute(dto *GetAllCategoryTaskDTO) ([]entities.CategoryTask, error) {
	categories, err := uc.CategoryTaskRepository.GetAll(dto.UserID)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
