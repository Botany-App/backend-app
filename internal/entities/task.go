package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	TaskDate     time.Time `json:"task_date"`
	UrgencyLevel int       `json:"urgency_level"`
	TaskStatus   string    `json:"task_status"`
	UserId       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CategoriesId []string  `json:"categories_id"`
	GardensId    []string  `json:"gardens_id"`
	PlantsId     []string  `json:"plants_id"`
}

type TaskOutputDTO struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	TaskDate     time.Time      `json:"task_date"`
	UrgencyLevel int            `json:"urgency_level"`
	TaskStatus   string         `json:"task_status"`
	UserId       string         `json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Categories   []CategoryTask `json:"categories"`
	Gardens      []Garden       `json:"gardens"`
	Plants       []Plant        `json:"plants"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) (string, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, userId, id string) error
	FindByID(ctx context.Context, userId, id string) (*TaskOutputDTO, error)
	FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*TaskOutputDTO, error)
	FindByName(ctx context.Context, userId, name string) ([]*TaskOutputDTO, error)
	FindAll(ctx context.Context, userId string) ([]*TaskOutputDTO, error)
	FindByStatus(ctx context.Context, userId, status string) ([]*TaskOutputDTO, error)
	FindByUrgencyLevel(ctx context.Context, userId string, urgencyLevel int) ([]*TaskOutputDTO, error)
}

func NewTask(
	name string,
	description string,
	taskDate time.Time,
	urgencyLevel int,
	taskStatus string,
	userId string,
	categoriesId []string,
	gardensId []string,
	plantsId []string,
) (*Task, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	if taskDate.IsZero() {
		return nil, errors.New("task date is required")
	}
	if urgencyLevel == 0 {
		return nil, errors.New("urgency level is required")
	}
	if taskStatus == "" {
		return nil, errors.New("task status is required")
	}
	if userId == "" {
		return nil, errors.New("user id is required")
	}

	return &Task{
		Id:           uuid.NewString(),
		Name:         name,
		Description:  description,
		TaskDate:     taskDate,
		UrgencyLevel: urgencyLevel,
		TaskStatus:   taskStatus,
		UserId:       userId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CategoriesId: categoriesId,
		GardensId:    gardensId,
		PlantsId:     plantsId,
	}, nil
}
