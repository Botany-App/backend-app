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
	UserID         uuid.UUID      `json:"user_id"`
	GardenPlantaId []uuid.UUID    `json:"garden_planta_id"`
	CategoriesId   []uuid.UUID    `json:"categories_id"`
	TaskStatus     TaskStatusEnum `json:"task_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type TaskRepository interface {
	Create(task *Task) error
	FindAll(userID uuid.UUID) ([]Task, error)
	FindByID(userID, id uuid.UUID) (*Task, error)
	FindAllByName(userID uuid.UUID, name string) ([]Task, error)
	FindAllByDate(userID uuid.UUID, date time.Time) ([]Task, error)
	FindAllByStatus(userID uuid.UUID, status TaskStatusEnum) ([]Task, error)
	FindTasksNearDeadline(userID uuid.UUID, days int) ([]Task, error)
	FindTasksFarFromDeadline(userID uuid.UUID, days int) ([]Task, error)
	Update(task *Task) error
	Delete(id uuid.UUID) error
}

func NewTask(name, description string, taskDate time.Time, userID uuid.UUID, gardenPlantaId, categoriesId []uuid.UUID, taskStatus TaskStatusEnum) (*Task, error) {
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
	if userID == uuid.Nil {
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
