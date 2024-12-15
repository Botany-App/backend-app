package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_task "github.com/lucasBiazon/botany-back/internal/usecases/task"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type TaskHandler struct {
	CreateTaskUseCase             *usecases_task.CreateTaskUseCase
	DeleteTaskUseCase             *usecases_task.DeleteTaskUseCase
	FindAllTaskUseCase            *usecases_task.FindAllTaskUseCase
	FindByCategoryNameTaskUseCase *usecases_task.FindByCategoryNameTaskUseCase
	FindByIdTaskUseCase           *usecases_task.FindByIdTaskUseCase
	UpdateTaskUseCase             *usecases_task.UpdateTaskUseCase
	FindByNameTaskUseCase         *usecases_task.FindByNameTaskUseCase
	FindByStatusTaskUseCase       *usecases_task.FindByStatusTaskUseCase
	FindByUrgencyLevelTaskUseCase *usecases_task.FindByUrgencyLevelTaskUseCase
}

func NewTaskHandler(
	createTaskUseCase *usecases_task.CreateTaskUseCase,
	deleteTaskUseCase *usecases_task.DeleteTaskUseCase,
	findAllTaskUseCase *usecases_task.FindAllTaskUseCase,
	findByCategoryNameTaskUseCase *usecases_task.FindByCategoryNameTaskUseCase,
	findByIdTaskUseCase *usecases_task.FindByIdTaskUseCase,
	updateTaskUseCase *usecases_task.UpdateTaskUseCase,
	findByNameTaskUseCase *usecases_task.FindByNameTaskUseCase,
	findByStatusTaskUseCase *usecases_task.FindByStatusTaskUseCase,
	findByUrgencyLevelTaskUseCase *usecases_task.FindByUrgencyLevelTaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		CreateTaskUseCase:             createTaskUseCase,
		DeleteTaskUseCase:             deleteTaskUseCase,
		FindAllTaskUseCase:            findAllTaskUseCase,
		FindByCategoryNameTaskUseCase: findByCategoryNameTaskUseCase,
		FindByIdTaskUseCase:           findByIdTaskUseCase,
		UpdateTaskUseCase:             updateTaskUseCase,
		FindByNameTaskUseCase:         findByNameTaskUseCase,
		FindByStatusTaskUseCase:       findByStatusTaskUseCase,
		FindByUrgencyLevelTaskUseCase: findByUrgencyLevelTaskUseCase,
	}
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.CreateTaskUseCaseInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	task, err := h.CreateTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "success", "Tarefa criada com sucesso", task)
}

func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.DeleteTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	err = h.DeleteTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefa deletada com sucesso", nil)
}

func (h *TaskHandler) FindAllTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.FindAllTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	tasks, err := h.FindAllTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefas encontradas com sucesso", tasks)
}

func (h *TaskHandler) FindByCategoryNameTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.FindByCategoryNameTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	tasks, err := h.FindByCategoryNameTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefas encontradas com sucesso", tasks)
}

func (h *TaskHandler) FindByIdTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.FindByIdTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	task, err := h.FindByIdTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefa encontrada com sucesso", task)
}

func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.UpdateTaskUseCaseInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	task, err := h.UpdateTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefa atualizada com sucesso", task)
}

func (h *TaskHandler) FindByNameTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.FindByNameTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	tasks, err := h.FindByNameTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefas encontradas com sucesso", tasks)
}

func (h *TaskHandler) FindByStatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.FindByStatusTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	tasks, err := h.FindByStatusTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefas encontradas com sucesso", tasks)
}

func (h *TaskHandler) FindByUrgencyLevelTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_task.FindByUrgencyLevelTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	tasks, err := h.FindByUrgencyLevelTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Tarefas encontradas com sucesso", tasks)
}

