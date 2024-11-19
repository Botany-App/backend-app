package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_categoryplant "github.com/lucasBiazon/botany-back/internal/usecases/category_plant"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type CategoryPlantHandler struct {
	CreateCategoryPlantUseCase     *usecases_categoryplant.CreateCategoryPlantUseCase
	FindAllCategoryPlantUseCase    *usecases_categoryplant.FindAllCategoryPlantUseCase
	FindByIdCategoryPlantUseCase   *usecases_categoryplant.FindByIdCategoryPlantUseCase
	FindByNameCategoryPlantUseCase *usecases_categoryplant.FindByNameCategoryPlantUseCase
	UpdateCategoryPlantUseCase     *usecases_categoryplant.UpdateCategoryPlantUseCase
	DeleteCategoryPlantUseCase     *usecases_categoryplant.DeleteCategoryPlantUseCase
}

func NewCategoryPlantHandler( // Construtor
	createCategoryPlantUseCase *usecases_categoryplant.CreateCategoryPlantUseCase,
	findAllCategoryPlantUseCase *usecases_categoryplant.FindAllCategoryPlantUseCase,
	findByIdCategoryPlantUseCase *usecases_categoryplant.FindByIdCategoryPlantUseCase,
	findByNameCategoryPlantUseCase *usecases_categoryplant.FindByNameCategoryPlantUseCase,
	updateCategoryPlantUseCase *usecases_categoryplant.UpdateCategoryPlantUseCase,
	deleteCategoryPlantUseCase *usecases_categoryplant.DeleteCategoryPlantUseCase,
) *CategoryPlantHandler {
	return &CategoryPlantHandler{
		CreateCategoryPlantUseCase:     createCategoryPlantUseCase,
		FindAllCategoryPlantUseCase:    findAllCategoryPlantUseCase,
		FindByIdCategoryPlantUseCase:   findByIdCategoryPlantUseCase,
		FindByNameCategoryPlantUseCase: findByNameCategoryPlantUseCase,
		UpdateCategoryPlantUseCase:     updateCategoryPlantUseCase,
		DeleteCategoryPlantUseCase:     deleteCategoryPlantUseCase,
	}
}

func (h *CategoryPlantHandler) CreateCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryPlant usecases_categoryplant.CreateCategoryPlantInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID

	if err := h.CreateCategoryPlantUseCase.Execute(context.Background(), categoryPlant); err != nil {
		log.Print(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "success", "Categoria de planta criada com sucesso", nil)
}

func (h *CategoryPlantHandler) FindAllCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	limit := utils.ParseQueryInt(r.URL.Query().Get("limit"), 10) // Default limit = 10
	offset := utils.ParseQueryInt(r.URL.Query().Get("offset"), 0)

	var categoryPlant usecases_categoryplant.FindAllCategoryPlantInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID
	categoryPlant.LIMIT = limit
	categoryPlant.OFFSET = offset
	categoryPlants, err := h.FindAllCategoryPlantUseCase.Execute(context.Background(), categoryPlant)
	if err != nil {
		log.Print(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categorias de plantas encontradas", categoryPlants)
}

func (h *CategoryPlantHandler) FindByIdCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	limit := utils.ParseQueryInt(r.URL.Query().Get("limit"), 10) // Default limit = 10
	offset := utils.ParseQueryInt(r.URL.Query().Get("offset"), 0)
	var categoryPlant usecases_categoryplant.FindByIdCategoryPlantInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID
	categoryPlant.LIMIT = limit
	categoryPlant.OFFSET = offset
	categoryPlantFound, err := h.FindByIdCategoryPlantUseCase.Execute(context.Background(), categoryPlant)
	if err != nil {
		log.Print(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta encontrada", categoryPlantFound)
}

func (h *CategoryPlantHandler) FindByNameCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	limit := utils.ParseQueryInt(r.URL.Query().Get("limit"), 10) // Default limit = 10
	offset := utils.ParseQueryInt(r.URL.Query().Get("offset"), 0)
	var categoryPlant usecases_categoryplant.FindByNameCategoryPlantInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.LIMIT = limit
	categoryPlant.OFFSET = offset
	categoryPlant.UserID = userID
	categoryPlantFound, err := h.FindByNameCategoryPlantUseCase.Execute(context.Background(), categoryPlant)
	if err != nil {
		log.Print(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta encontrada", categoryPlantFound)
}

func (h *CategoryPlantHandler) UpdateCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryPlant usecases_categoryplant.UpdateCategoryPlantInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID

	if err := h.UpdateCategoryPlantUseCase.Execute(context.Background(), categoryPlant); err != nil {
		log.Print(err)
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta atualizada com sucesso", nil)
}

func (h *CategoryPlantHandler) DeleteCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryPlant usecases_categoryplant.DeleteCategoryPlantInputDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryPlant.UserID = userID

	if err := h.DeleteCategoryPlantUseCase.Execute(context.Background(), categoryPlant); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	utils.JsonResponse(w, http.StatusOK, "success", "Categoria de planta deletada com sucesso", nil)
}
