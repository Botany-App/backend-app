package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskStatusEnum string

const (
	Pending    TaskStatusEnum = "pending"
	InProgress TaskStatusEnum = "in_progress"
	Completed  TaskStatusEnum = "completed"
)

type Task struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	TaskDate       time.Time      `json:"task_date"`
	UserID         string         `json:"user_id"`
	GardenPlantaId []string       `json:"garden_planta_id"`
	CategoriesId   []string       `json:"categories_id"`
	TaskStatus     TaskStatusEnum `json:"task_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type TaskRepository interface {
	Create(task *Task) error
	AddCategory(taskID, categoryID string) error
	AddGardenPlanta(taskID, gardenPlantaID string) error
	FindAll(userID string) ([]Task, error)
	FindByID(userID, id string) (*Task, error)
	FindAllByName(userID string, name string) ([]Task, error)
	FindAllByDate(userID string, date time.Time) ([]Task, error)
	FindAllByStatus(userID string, status TaskStatusEnum) ([]Task, error)
	FindTasksNearDeadline(userID string, days int) ([]Task, error)
	FindTasksFarFromDeadline(userID string, days int) ([]Task, error)
	FindAllByCategory(userID string, categoryID string) ([]Task, error)
	Update(task *Task) error
	UpdateTaskCategory(taskID, categoryID string) error
	UpdateTaskGardenPlanta(taskID, gardenPlantaID string) error
	Delete(userID string, id string) error
}

func NewTask(name, description string, taskDate time.Time, userID string, gardenPlantaId, categoriesId []string, taskStatus TaskStatusEnum) (*Task, error) {
	if taskStatus == "" {
		taskStatus = Pending
	}
	if len(gardenPlantaId) == 0 {
		gardenPlantaId = nil
	}
	if len(categoriesId) == 0 {
		categoriesId = nil
	}
	if taskDate.IsZero() {
		taskDate = time.Now()
	}
	if description == "" {
		description = "No description"
	}

	if name == "" {
		return nil, errors.New("task name is required")
	}
	if len(name) > 50 {
		return nil, errors.New("task name exceeds 50 characters")
	}
	if len(description) > 100 {
		return nil, errors.New("task description exceeds 100 characters")
	}
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if taskDate.Before(time.Now()) {
		return nil, errors.New("task date cannot be in the past")
	}
	return &Task{
		ID:             uuid.New(),
		Name:           name,
		Description:    description,
		TaskDate:       taskDate,
		UserID:         userID,
		GardenPlantaId: gardenPlantaId,
		CategoriesId:   categoriesId,
		TaskStatus:     taskStatus,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}
