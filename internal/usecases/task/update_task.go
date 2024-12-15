package usecases_task

import (
	"context"
	"errors"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateTaskUseCase struct {
	Repository entities.TaskRepository
}

type UpdateTaskUseCaseInputDTO struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TaskDate     time.Time `json:"task_date"`
	UrgencyLevel int       `json:"urgency_level"`
	TaskStatus   string    `json:"task_status"`
	UserId       string    `json:"user_id"`
	CategoriesId []string  `json:"categories_id"`
	GardensId    []string  `json:"gardens_id"`
	PlantsId     []string  `json:"plants_id"`
}

func NewUpdateTaskUseCase(repository entities.TaskRepository) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		Repository: repository,
	}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, input UpdateTaskUseCaseInputDTO) (*entities.TaskOutputDTO, error) {
	// Verifica se a tarefa existe pelo ID
	existingTask, err := uc.Repository.FindByID(ctx, input.UserId, input.Id)
	if err != nil {
		return nil, err
	}
	if existingTask == nil {
		return nil, errors.New("tarefa não encontrada")
	}

	// Verifica se já existe uma tarefa com o mesmo nome para o usuário
	taskByName, err := uc.Repository.FindByName(ctx, input.UserId, input.Name)
	if err != nil {
		return nil, err
	}
	for _, task := range taskByName {
		if task.Id != input.Id {
			return nil, errors.New("já existe uma tarefa com este nome")
		}
	}

	// Cria um mapa para os campos atualizados
	updatedFields := make(map[string]interface{})

	if existingTask.Name != input.Name {
		updatedFields["name"] = input.Name
	}
	if existingTask.Description != input.Description {
		updatedFields["description"] = input.Description
	}
	if !existingTask.TaskDate.Equal(input.TaskDate) {
		updatedFields["task_date"] = input.TaskDate
	}
	if existingTask.UrgencyLevel != input.UrgencyLevel {
		updatedFields["urgency_level"] = input.UrgencyLevel
	}
	if existingTask.TaskStatus != input.TaskStatus {
		updatedFields["task_status"] = input.TaskStatus
	}

	// Verifica e atualiza as categorias
	existingCategories := extractCategoryIDs(existingTask.Categories)
	if len(existingCategories) != len(input.CategoriesId) || !equalStringSlices(existingCategories, input.CategoriesId) {
		updatedFields["categories_id"] = input.CategoriesId
	}

	// Verifica e atualiza os jardins
	existingGardens := extractGardenIDs(existingTask.Gardens)
	if len(existingGardens) != len(input.GardensId) || !equalStringSlices(existingGardens, input.GardensId) {
		updatedFields["gardens_id"] = input.GardensId
	}

	// Verifica e atualiza as plantas
	existingPlants := extractPlantIDs(existingTask.Plants)
	if len(existingPlants) != len(input.PlantsId) || !equalStringSlices(existingPlants, input.PlantsId) {
		updatedFields["plants_id"] = input.PlantsId
	}

	// Cria a tarefa atualizada
	updatedTask := &entities.Task{
		Id:           input.Id,
		Name:         input.Name,
		Description:  input.Description,
		TaskDate:     input.TaskDate,
		UrgencyLevel: input.UrgencyLevel,
		TaskStatus:   input.TaskStatus,
		UserId:       input.UserId,
		CategoriesId: input.CategoriesId,
		GardensId:    input.GardensId,
		PlantsId:     input.PlantsId,
		UpdatedAt:    time.Now(),
	}

	// Atualiza a tarefa no repositório
	if err := uc.Repository.Update(ctx, updatedTask); err != nil {
		return nil, err
	}

	// Retorna a tarefa atualizada
	updatedTaskOutput, err := uc.Repository.FindByID(ctx, input.UserId, input.Id)
	if err != nil {
		return nil, err
	}
	return updatedTaskOutput, nil
}

// Funções auxiliares
func extractCategoryIDs(categories []entities.CategoryTask) []string {
	ids := make([]string, len(categories))
	for i, category := range categories {
		ids[i] = category.Id
	}
	return ids
}

func extractGardenIDs(gardens []entities.Garden) []string {
	ids := make([]string, len(gardens))
	for i, garden := range gardens {
		ids[i] = garden.Id
	}
	return ids
}

func extractPlantIDs(plants []entities.Plant) []string {
	ids := make([]string, len(plants))
	for i, plant := range plants {
		ids[i] = plant.Id
	}
	return ids
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)
	for _, v := range a {
		m[v]++
	}
	for _, v := range b {
		if m[v] == 0 {
			return false
		}
		m[v]--
	}
	return true
}
