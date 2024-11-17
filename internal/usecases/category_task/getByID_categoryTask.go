package usecases_categorytask

import "github.com/lucasBiazon/botany-back/internal/entities"

type GetByIdCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type GetByIdCategoryTaskDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func NewGetByIdCategoryTaskUseCase(categoryTaskRepository entities.CategoryTaskRepository) *GetByIdCategoryTaskUseCase {
	return &GetByIdCategoryTaskUseCase{
		CategoryTaskRepository: categoryTaskRepository,
	}
}
func (uc *GetByIdCategoryTaskUseCase) Execute(dto *GetByIdCategoryTaskDTO) (*entities.CategoryTask, error) {
	categoriesTask, err := uc.CategoryTaskRepository.GetByID(dto.ID, dto.UserID)
	if err != nil {
		return nil, err
	}
	categoryTask := categoriesTask[0]
	return &categoryTask, nil
}
