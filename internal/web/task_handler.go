package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/lucasBiazon/botany-back/internal/entities"
	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/tasks"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type TaskHandlers struct {
	CreateTaskUseCase               *usecases.CreateTaskUseCase
	DeleteTaskUseCase               *usecases.DeleteTaskUseCase
	FindAllByCategoryTaskUseCase    *usecases.FindAllByCategoryTaskUseCase
	FindAllByDateTaskUseCase        *usecases.FindAllByDateTaskUseCase
	FindAllByNameTaskUseCase        *usecases.FindAllByNameTaskUseCase
	FindAllByStatusTaskUseCase      *usecases.FindAllByStatusTaskUseCase
	FindAllTaskUseCase              *usecases.FindAllTaskUseCase
	FindByIDTaskUseCase             *usecases.FindByIDTaskUseCase
	UpdateTaskUseCase               *usecases.UpdateTaskUseCase
	FindTaskNearDeadLineTaskUseCase *usecases.FindTaskNearDeadLineTaskUseCase
	FindTasksFarFromDeadlineUseCase *usecases.FindTasksFarFromDeadlineUseCase
}

func NewTaskHandlers(createTaskUseCase *usecases.CreateTaskUseCase, deleteTaskUseCase *usecases.DeleteTaskUseCase,
	findAllByCategoryTaskUseCase *usecases.FindAllByCategoryTaskUseCase, findAllByDateTaskUseCase *usecases.FindAllByDateTaskUseCase,
	findAllByNameTaskUseCase *usecases.FindAllByNameTaskUseCase, findAllByStatusTaskUseCase *usecases.FindAllByStatusTaskUseCase,
	findAllTaskUseCase *usecases.FindAllTaskUseCase, findByIDTaskUseCase *usecases.FindByIDTaskUseCase, updateTaskUseCase *usecases.UpdateTaskUseCase,
	findTaskNearDeadLineTaskUseCase *usecases.FindTaskNearDeadLineTaskUseCase,
	findTasksFarFromDeadlineUseCase *usecases.FindTasksFarFromDeadlineUseCase) *TaskHandlers {
	return &TaskHandlers{
		CreateTaskUseCase:               createTaskUseCase,
		DeleteTaskUseCase:               deleteTaskUseCase,
		FindAllByCategoryTaskUseCase:    findAllByCategoryTaskUseCase,
		FindAllByDateTaskUseCase:        findAllByDateTaskUseCase,
		FindAllByNameTaskUseCase:        findAllByNameTaskUseCase,
		FindAllByStatusTaskUseCase:      findAllByStatusTaskUseCase,
		FindAllTaskUseCase:              findAllTaskUseCase,
		FindByIDTaskUseCase:             findByIDTaskUseCase,
		UpdateTaskUseCase:               updateTaskUseCase,
		FindTaskNearDeadLineTaskUseCase: findTaskNearDeadLineTaskUseCase,
		FindTasksFarFromDeadlineUseCase: findTasksFarFromDeadlineUseCase,
	}
}

// CreateTaskHandler is a handler for creating a task
func (h *TaskHandlers) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var input usecases.CreateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Error decoding request", nil)
		return
	}
	input.UserID = userID
	if err := h.CreateTaskUseCase.Execute(context.Background(), &input); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error creating task", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Task created", nil)
}

// DeleteTaskHandler is a handler for deleting a task
func (h *TaskHandlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	var input usecases.DeleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Error decoding request", nil)
		return
	}
	input.UserID = userID
	if err := h.DeleteTaskUseCase.Execute(context.Background(), &input); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error deleting task", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Task deleted", nil)
}

// FindAllByCategoryTaskHandler is a handler for finding all tasks by category
func (h *TaskHandlers) FindAllByCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindAllByCategoryDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Category not found", nil)
		return
	}
	input.CategoryID = category
	input.UserID = userID
	tasks, err := h.FindAllByCategoryTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}

// FindAllByDateTaskHandler is a handler for finding all tasks by date
func (h *TaskHandlers) FindAllByDateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindAllByDateDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Date not found", nil)
		return
	}
	input.Date = date
	input.UserID = userID
	tasks, err := h.FindAllByDateTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}

// FindAllByNameTaskHandler is a handler for finding all tasks by name
func (h *TaskHandlers) FindAllByNameTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindAllByNameDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Name not found", nil)
		return
	}
	input.Name = name
	input.UserID = userID
	tasks, err := h.FindAllByNameTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}

// FindAllByStatusTaskHandler is a handler for finding all tasks by status
func (h *TaskHandlers) FindAllByStatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindAllByStatusDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Status not found", nil)
		return
	}
	input.Status = entities.TaskStatusEnum(status)
	input.UserID = userID
	tasks, err := h.FindAllByStatusTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}

// FindAllTaskHandler is a handler for finding all tasks
func (h *TaskHandlers) FindAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	var input usecases.FindAllTaskDTO
	input.UserID = userID
	tasks, err := h.FindAllTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}

// FindByIDTaskHandler is a handler for finding a task by ID
func (h *TaskHandlers) FindByIDTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindByIDTaskDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "ID not found", nil)
		return
	}
	input.ID = id
	input.UserID = userID
	task, err := h.FindByIDTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding task", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Task found", task)
}

// UpdateTaskHandler is a handler for updating a task
func (h *TaskHandlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.UpdateTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Error decoding request", nil)
		return
	}
	if err := h.UpdateTaskUseCase.Execute(context.Background(), &input); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error updating task", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Task updated", nil)
}

// FindTaskNearDeadLineTaskHandler is a handler for finding tasks near deadline
func (h *TaskHandlers) FindTaskNearDeadLineTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindTasksNearDeadlineDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	input.UserID = userID
	tasks, err := h.FindTaskNearDeadLineTaskUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}

// FindTasksFarFromDeadlineHandler is a handler for finding tasks far from deadline
func (h *TaskHandlers) FindTasksFarFromDeadlineHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.FindTasksFarFromDeadlineDTO
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Invalid or expired token", nil)
		return
	}

	input.UserID = userID
	tasks, err := h.FindTasksFarFromDeadlineUseCase.Execute(context.Background(), &input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Error finding tasks", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tasks found", tasks)
}
