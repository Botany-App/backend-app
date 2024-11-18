package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_categoryplant "github.com/lucasBiazon/botany-back/internal/usecases/category_plant"
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
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryPlant usecases_categoryplant.CreateCategoryPlantDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID

	if err := h.CreateCategoryPlantUseCase.Execute(context.Background(), categoryPlant); err != nil {
		log.Print(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
}

func (h *CategoryPlantHandler) FindAllCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryPlant usecases_categoryplant.FindAllCategoryPlantDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID
	categoryPlants, err := h.FindAllCategoryPlantUseCase.Execute(context.Background(), categoryPlant)
	if err != nil {
		log.Print(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categorias de plantas encontradas", categoryPlants)
}

func (h *CategoryPlantHandler) FindByIdCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryPlant usecases_categoryplant.FindByIdCategoryPlantDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID
	categoryPlantFound, err := h.FindByIdCategoryPlantUseCase.Execute(context.Background(), categoryPlant)
	if err != nil {
		log.Print(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de planta encontrada", categoryPlantFound)
}

func (h *CategoryPlantHandler) FindByNameCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryPlant usecases_categoryplant.FindByNameCategoryPlantDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID
	categoryPlantFound, err := h.FindByNameCategoryPlantUseCase.Execute(context.Background(), categoryPlant)
	if err != nil {
		log.Print(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de planta encontrada", categoryPlantFound)
}

func (h *CategoryPlantHandler) UpdateCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}

	var categoryPlant usecases_categoryplant.UpdateCategoryPlantDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		log.Println(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	categoryPlant.UserID = userID

	if err := h.UpdateCategoryPlantUseCase.Execute(context.Background(), categoryPlant); err != nil {
		log.Print(err)
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de planta atualizada com sucesso", nil)
}

func (h *CategoryPlantHandler) DeleteCategoryPlantHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var categoryPlant usecases_categoryplant.DeleteCategoryPlantDTO
	if err := json.NewDecoder(r.Body).Decode(&categoryPlant); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	categoryPlant.UserID = userID

	if err := h.DeleteCategoryPlantUseCase.Execute(context.Background(), categoryPlant); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	jsonResponse(w, http.StatusOK, "success", "Categoria de planta deletada com sucesso", nil)
}
