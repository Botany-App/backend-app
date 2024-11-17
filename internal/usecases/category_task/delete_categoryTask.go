package usecases_categorytask

import "github.com/lucasBiazon/botany-back/internal/entities"

type DeleteCategoryTaskUseCase struct {
	CategoryTaskRepository entities.CategoryTaskRepository
}

type DeleteCategoryTaskDTO struct {
	UserID string `json:"user_id"`
	ID     string `json:"id"`
}

func NewDeleteCategoryTaskUseCase(repository entities.CategoryTaskRepository) *DeleteCategoryTaskUseCase {
	return &DeleteCategoryTaskUseCase{CategoryTaskRepository: repository}
}

func (useCase DeleteCategoryTaskUseCase) Execute(dto DeleteCategoryTaskDTO) error {
	err := useCase.CategoryTaskRepository.Delete(dto.UserID, dto.ID)
	if err != nil {
		return err
	}

	return nil
}
