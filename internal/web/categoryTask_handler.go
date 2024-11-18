package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_categorytask "github.com/lucasBiazon/botany-back/internal/usecases/category_task"
)

type CategoryTaskHandler struct {
	CreateCategoryTaskUseCase    *usecases_categorytask.CreateCategoryTaskUseCase
	DeleteCategoryTaskUseCase    *usecases_categorytask.DeleteCategoryTaskUseCase
	GetByNameCategoryTaskUseCase *usecases_categorytask.GetByNameCategoryTaskUseCase
	UpdateCategoryTaskUseCase    *usecases_categorytask.UpdateCategoryTaskUseCase
	GetByIdCategoryTaskUseCase   *usecases_categorytask.GetByIdCategoryTaskUseCase
	GetAllCategoryTaskUseCase    *usecases_categorytask.GetAllCategoryTaskUseCase
}

func NewCategoryTaskHandler(
	createCategoryTaskUseCase *usecases_categorytask.CreateCategoryTaskUseCase,
	deleteCategoryTaskUseCase *usecases_categorytask.DeleteCategoryTaskUseCase,
	getByNameCategoryTaskUseCase *usecases_categorytask.GetByNameCategoryTaskUseCase,
	updateCategoryTaskUseCase *usecases_categorytask.UpdateCategoryTaskUseCase,
	getByIdCategoryTaskUseCase *usecases_categorytask.GetByIdCategoryTaskUseCase,
	getAllCategoryTaskUseCase *usecases_categorytask.GetAllCategoryTaskUseCase,
) *CategoryTaskHandler {
	return &CategoryTaskHandler{
		CreateCategoryTaskUseCase:    createCategoryTaskUseCase,
		DeleteCategoryTaskUseCase:    deleteCategoryTaskUseCase,
		GetByNameCategoryTaskUseCase: getByNameCategoryTaskUseCase,
		UpdateCategoryTaskUseCase:    updateCategoryTaskUseCase,
		GetByIdCategoryTaskUseCase:   getByIdCategoryTaskUseCase,
		GetAllCategoryTaskUseCase:    getAllCategoryTaskUseCase,
	}
}

func (h *CategoryTaskHandler) CreateCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryTask usecases_categorytask.CreateCategoryTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryTask.UserID = userID

	if err := h.CreateCategoryTaskUseCase.Execute(context.Background(), categoryTask); err != nil {
		log.Print(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	jsonResponse(w, http.StatusCreated, "success", "Categoria de tarefa criada com sucesso", nil)
}

func (h *CategoryTaskHandler) DeleteCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryTask usecases_categorytask.DeleteCategoryTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID

	if err := h.DeleteCategoryTaskUseCase.Execute(context.Background(), categoryTask); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de tarefa deletada com sucesso", nil)
}

func (h *CategoryTaskHandler) GetByNameCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryTask usecases_categorytask.GetByNameCategoryTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID

	categoryTasks, err := h.GetByNameCategoryTaskUseCase.Execute(context.Background(), &categoryTask)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de tarefa encontrada", categoryTasks)
}

func (h *CategoryTaskHandler) UpdateCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryTask usecases_categorytask.UpdateCategoryTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID

	if err := h.UpdateCategoryTaskUseCase.Execute(context.Background(), categoryTask); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de tarefa atualizada com sucesso", nil)
}

func (h *CategoryTaskHandler) GetByIdCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryTask usecases_categorytask.GetByIdCategoryTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID

	categoryTaskFound, err := h.GetByIdCategoryTaskUseCase.Execute(context.Background(), &categoryTask)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de tarefa encontrada", categoryTaskFound)
}

func (h *CategoryTaskHandler) GetAllCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryTask usecases_categorytask.GetAllCategoryTaskDTO
	categoryTask.UserID = userID
	categoryTasks, err := h.GetAllCategoryTaskUseCase.Execute(context.Background(), &categoryTask)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Categorias de tarefas encontradas", categoryTasks)
}
