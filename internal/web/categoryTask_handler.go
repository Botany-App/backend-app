package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_categorytask "github.com/lucasBiazon/botany-back/internal/usecases/category_task"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type CategoryTaskHandler struct {
	CreateCategoryTaskUseCase     *usecases_categorytask.CreateCategoryTaskUseCase
	DeleteCategoryTaskUseCase     *usecases_categorytask.DeleteCategoryTaskUseCase
	FindByNameCategoryTaskUseCase *usecases_categorytask.FindByNameCategoryTaskUseCase
	UpdateCategoryTaskUseCase     *usecases_categorytask.UpdateCategoryTaskUseCase
	FindByIdCategoryTaskUseCase   *usecases_categorytask.FindByIdCategoryTaskUseCase
	FindAllCategoryTaskUseCase    *usecases_categorytask.FindAllCategoryTaskUseCase
}

func NewCategoryTaskHandler(
	createCategoryTaskUseCase *usecases_categorytask.CreateCategoryTaskUseCase,
	deleteCategoryTaskUseCase *usecases_categorytask.DeleteCategoryTaskUseCase,
	findByNameCategoryTaskUseCase *usecases_categorytask.FindByNameCategoryTaskUseCase,
	updateCategoryTaskUseCase *usecases_categorytask.UpdateCategoryTaskUseCase,
	findByIdCategoryTaskUseCase *usecases_categorytask.FindByIdCategoryTaskUseCase,
	findAllCategoryTaskUseCase *usecases_categorytask.FindAllCategoryTaskUseCase,
) *CategoryTaskHandler {
	return &CategoryTaskHandler{
		CreateCategoryTaskUseCase:     createCategoryTaskUseCase,
		DeleteCategoryTaskUseCase:     deleteCategoryTaskUseCase,
		FindByNameCategoryTaskUseCase: findByNameCategoryTaskUseCase,
		UpdateCategoryTaskUseCase:     updateCategoryTaskUseCase,
		FindByIdCategoryTaskUseCase:   findByIdCategoryTaskUseCase,
		FindAllCategoryTaskUseCase:    findAllCategoryTaskUseCase,
	}
}

func (h *CategoryTaskHandler) CreateCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryTask usecases_categorytask.CreateCategoryTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		log.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryTask.UserID = userID

	if err := h.CreateCategoryTaskUseCase.Execute(context.Background(), categoryTask); err != nil {
		log.Print(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "success", "Categoria de tarefa criada com sucesso", nil)
}

func (h *CategoryTaskHandler) DeleteCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryTask usecases_categorytask.DeleteCategoryTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID

	if err := h.DeleteCategoryTaskUseCase.Execute(context.Background(), categoryTask); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de tarefa deletada com sucesso", nil)
}

func (h *CategoryTaskHandler) FindByNameCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	limit := utils.ParseQueryInt(r.URL.Query().Get("limit"), 10) // Default limit = 10
	offset := utils.ParseQueryInt(r.URL.Query().Get("offset"), 0)
	var categoryTask usecases_categorytask.FindByNameCategoryTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID
	categoryTask.LIMIT = limit
	categoryTask.OFFSET = offset
	categoryTasks, err := h.FindByNameCategoryTaskUseCase.Execute(context.Background(), categoryTask)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de tarefa encontrada", categoryTasks)
}

func (h *CategoryTaskHandler) UpdateCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryTask usecases_categorytask.UpdateCategoryTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID

	if err := h.UpdateCategoryTaskUseCase.Execute(context.Background(), categoryTask); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de tarefa atualizada com sucesso", nil)
}

func (h *CategoryTaskHandler) FindByIdCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	limit := utils.ParseQueryInt(r.URL.Query().Get("limit"), 10) // Default limit = 10
	offset := utils.ParseQueryInt(r.URL.Query().Get("offset"), 0)
	var categoryTask usecases_categorytask.FindByIdCategoryTaskInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryTask); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryTask.UserID = userID
	categoryTask.LIMIT = limit
	categoryTask.OFFSET = offset
	categoryTaskFound, err := h.FindByIdCategoryTaskUseCase.Execute(context.Background(), categoryTask)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de tarefa encontrada", categoryTaskFound)
}

func (h *CategoryTaskHandler) FindAllCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	limit := utils.ParseQueryInt(r.URL.Query().Get("limit"), 10) // Default limit = 10
	offset := utils.ParseQueryInt(r.URL.Query().Get("offset"), 0)
	var categoryTask usecases_categorytask.FindAllCategoryTaskInputDTO
	categoryTask.UserID = userID
	categoryTask.LIMIT = limit
	categoryTask.OFFSET = offset
	categoryTasks, err := h.FindAllCategoryTaskUseCase.Execute(context.Background(), categoryTask)
	if err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categorias de tarefas encontradas", categoryTasks)
}
