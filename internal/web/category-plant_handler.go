package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecase_categoryplant "github.com/lucasBiazon/botany-back/internal/usecases/category-plant"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type CategoryPlantHandlers struct {
	CreateCategoryPlantUseCase     *usecase_categoryplant.CreateCategoryPlantUseCase
	DeleteCategoryPlantUseCase     *usecase_categoryplant.DeleteCategoryPlantUseCase
	FindAllCategoryPlantUseCase    *usecase_categoryplant.FindAllCategoryPlantUseCase
	FindByIdCategoryPlantUseCase   *usecase_categoryplant.FindByIdCategoryPlantUseCase
	FindByNameCategoryPlantUseCase *usecase_categoryplant.FindByNameCategoryPlantUseCase
	UpdateCategoryPlantUseCase     *usecase_categoryplant.UpdateCategoryPlantUseCase
}

func NewCategoryPlantHandler(
	createCategoryPlantUseCase *usecase_categoryplant.CreateCategoryPlantUseCase,
	deleteCategoryPlantUseCase *usecase_categoryplant.DeleteCategoryPlantUseCase,
	findAllCategoryPlantUseCase *usecase_categoryplant.FindAllCategoryPlantUseCase,
	findByIdCategoryPlantUseCase *usecase_categoryplant.FindByIdCategoryPlantUseCase,
	findByNameCategoryPlantUseCase *usecase_categoryplant.FindByNameCategoryPlantUseCase,
	updateCategoryPlantUseCase *usecase_categoryplant.UpdateCategoryPlantUseCase,
) *CategoryPlantHandlers {
	return &CategoryPlantHandlers{
		CreateCategoryPlantUseCase:     createCategoryPlantUseCase,
		DeleteCategoryPlantUseCase:     deleteCategoryPlantUseCase,
		FindAllCategoryPlantUseCase:    findAllCategoryPlantUseCase,
		FindByIdCategoryPlantUseCase:   findByIdCategoryPlantUseCase,
		FindByNameCategoryPlantUseCase: findByNameCategoryPlantUseCase,
		UpdateCategoryPlantUseCase:     updateCategoryPlantUseCase,
	}
}

func (h *CategoryPlantHandlers) CreateCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase_categoryplant.CreateCategoryPlantInputDTO
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
	categoryPlant, err := h.CreateCategoryPlantUseCase.Execute(context.Background(), input, userId)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao criar categoria de planta", err.Error())
		return
	}
	if categoryPlant == nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Category já existente", nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "sucess", "Categoria de planta criada com sucesso", categoryPlant)
}

func (h *CategoryPlantHandlers) DeleteCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase_categoryplant.DeleteCategoryPlantInputDTO
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
	if err := h.DeleteCategoryPlantUseCase.Execute(context.Background(), input); err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao deletar categoria de planta", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta deletado com sucesso", nil)
}

func (h *CategoryPlantHandlers) FindAllCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	categoryPlants, err := h.FindAllCategoryPlantUseCase.Execute(context.Background(), userId)
	if len(categoryPlants) == 0 {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Categorias de planta não encontradas", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar categorias de planta", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categorias de planta encontradas", categoryPlants)
}

func (h *CategoryPlantHandlers) FindByIdCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase_categoryplant.FindByIdCategoryPlantInputDTO
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
	categoryPlant, err := h.FindByIdCategoryPlantUseCase.Execute(context.Background(), input)
	if categoryPlant == nil {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Categoria de planta não encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar categoria de planta", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta encontrada", categoryPlant)

}

func (h *CategoryPlantHandlers) FindByNameCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase_categoryplant.FindByNameCategoryPlantInputDTO
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
	categoryPlant, err := h.FindByNameCategoryPlantUseCase.Execute(context.Background(), input)
	if len(categoryPlant) == 0 {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Categoria de planta não encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar categoria de planta", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta encontrada", categoryPlant)
}

func (h *CategoryPlantHandlers) UpdateCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase_categoryplant.UpdateCategoryPlantInputDTO
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
	categoryPlant, err := h.UpdateCategoryPlantUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao atualizar categoria de planta", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta atualizada com sucesso", categoryPlant)
}
