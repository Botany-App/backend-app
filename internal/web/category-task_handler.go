package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_categoryTask "github.com/lucasBiazon/botany-back/internal/usecases/category-task"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type CategoryTaskHandlers struct {
	CreateCategoryTaskUseCase     *usecases_categoryTask.CreateCategoryTaskUseCase
	DeleteCategoryTaskUseCase     *usecases_categoryTask.DeleteCategoryTaskUseCase
	FindAllCategoryTaskUseCase    *usecases_categoryTask.FindAllCategoryTaskUseCase
	FindByIdCategoryTaskUseCase   *usecases_categoryTask.FindByIdCategoryTaskUseCase
	FindByNameCategoryTaskUseCase *usecases_categoryTask.FindByNameCategoryTaskUseCase
	UpdateCategoryTaskUseCase     *usecases_categoryTask.UpdateCategoryTaskUseCase
}

func NewCategoryTaskHandler(
	createCategoryTaskUseCase *usecases_categoryTask.CreateCategoryTaskUseCase,
	deleteCategoryTaskUseCase *usecases_categoryTask.DeleteCategoryTaskUseCase,
	findAllCategoryTaskUseCase *usecases_categoryTask.FindAllCategoryTaskUseCase,
	findByIdCategoryTaskUseCase *usecases_categoryTask.FindByIdCategoryTaskUseCase,
	findByNameCategoryTaskUseCase *usecases_categoryTask.FindByNameCategoryTaskUseCase,
	updateCategoryTaskUseCase *usecases_categoryTask.UpdateCategoryTaskUseCase,
) *CategoryTaskHandlers {
	return &CategoryTaskHandlers{
		CreateCategoryTaskUseCase:     createCategoryTaskUseCase,
		DeleteCategoryTaskUseCase:     deleteCategoryTaskUseCase,
		FindAllCategoryTaskUseCase:    findAllCategoryTaskUseCase,
		FindByIdCategoryTaskUseCase:   findByIdCategoryTaskUseCase,
		FindByNameCategoryTaskUseCase: findByNameCategoryTaskUseCase,
		UpdateCategoryTaskUseCase:     updateCategoryTaskUseCase,
	}
}

func (h *CategoryTaskHandlers) CreateCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_categoryTask.CreateCategoryTaskInputDTO
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
	categoryTask, err := h.CreateCategoryTaskUseCase.Execute(context.Background(), input, userId)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao criar categoria de Taska", err.Error())
		return
	}
	if categoryTask == nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Category já existente", nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "sucess", "Categoria de Taska criada com sucesso", categoryTask)
}

func (h *CategoryTaskHandlers) DeleteCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_categoryTask.DeleteCategoryTaskInputDTO
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
	if err := h.DeleteCategoryTaskUseCase.Execute(context.Background(), input); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao deletar categoria de Taska", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de Taska deletado com sucesso", nil)
}

func (h *CategoryTaskHandlers) FindAllCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	categoryTasks, err := h.FindAllCategoryTaskUseCase.Execute(context.Background(), userId)
	if len(categoryTasks) == 0 {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Categorias de Taska não encontradas", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar categorias de Taska", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categorias de Taska encontradas", categoryTasks)
}

func (h *CategoryTaskHandlers) FindByIdCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_categoryTask.FindByIdCategoryTaskInputDTO
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
	categoryTask, err := h.FindByIdCategoryTaskUseCase.Execute(context.Background(), input)
	if categoryTask == nil {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Categoria de Taska não encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar categoria de Taska", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de Taska encontrada", categoryTask)

}

func (h *CategoryTaskHandlers) FindByNameCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_categoryTask.FindByNameCategoryTaskInputDTO
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
	categoryTask, err := h.FindByNameCategoryTaskUseCase.Execute(context.Background(), input)
	if len(categoryTask) == 0 {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Categoria de Taska não encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar categoria de Taska", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de Taska encontrada", categoryTask)
}

func (h *CategoryTaskHandlers) UpdateCategoryTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_categoryTask.UpdateCategoryTaskInputDTO
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
	categoryTask, err := h.UpdateCategoryTaskUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao atualizar categoria de Taska", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de Taska atualizada com sucesso", categoryTask)
}
