package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_plant "github.com/lucasBiazon/botany-back/internal/usecases/plant"
	usecases_specie "github.com/lucasBiazon/botany-back/internal/usecases/specie"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type PlantHandler struct {
	CreatePlantUseCase             *usecases_plant.CreatePlantUseCase
	DeletePlantUseCase             *usecases_plant.DeletePlantUseCase
	FindAllPlantUseCase            *usecases_plant.FindAllPlantUseCase
	FindByCategoryNamePlantUseCase *usecases_plant.FindByCategoryNamePlantUseCase
	FindByIdPlantUseCase           *usecases_plant.FindByIdPlantUseCase
	FindByNamePlantUseCase         *usecases_plant.FindByNamePlantUseCase
	FindBySpecieNamePlantUseCase   *usecases_plant.FindBySpecieNamePlantUseCase
	UpdatePlantUseCase             *usecases_plant.UpdatePlantUseCase
	FindByIdSpecieUseCase          *usecases_specie.FindByIdSpecieUseCase
}

func NewPlantHandler(
	createPlantUseCase *usecases_plant.CreatePlantUseCase,
	deletePlantUseCase *usecases_plant.DeletePlantUseCase,
	findAllPlantUseCase *usecases_plant.FindAllPlantUseCase,
	findByCategoryNamePlantUseCase *usecases_plant.FindByCategoryNamePlantUseCase,
	findByIdPlantUseCase *usecases_plant.FindByIdPlantUseCase,
	findByNamePlantUseCase *usecases_plant.FindByNamePlantUseCase,
	findBySpecieNamePlantUseCase *usecases_plant.FindBySpecieNamePlantUseCase,
	updatePlantUseCase *usecases_plant.UpdatePlantUseCase,
	findByIdSpecieUseCase *usecases_specie.FindByIdSpecieUseCase,
) *PlantHandler {
	return &PlantHandler{
		CreatePlantUseCase:             createPlantUseCase,
		DeletePlantUseCase:             deletePlantUseCase,
		FindAllPlantUseCase:            findAllPlantUseCase,
		FindByCategoryNamePlantUseCase: findByCategoryNamePlantUseCase,
		FindByIdPlantUseCase:           findByIdPlantUseCase,
		FindByNamePlantUseCase:         findByNamePlantUseCase,
		FindBySpecieNamePlantUseCase:   findBySpecieNamePlantUseCase,
		UpdatePlantUseCase:             updatePlantUseCase,
		FindByIdSpecieUseCase:          findByIdSpecieUseCase,
	}
}

func (h *PlantHandler) CreatePlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.CreatePlantUseCaseInputDTO
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
	log.Println("userId", input.SpeciesID)
	specieHaverstTime, err := h.FindByIdSpecieUseCase.Execute(r.Context(), usecases_specie.FindByIdSpecieInputDTO{Id: input.SpeciesID})
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	input.UserID = userId
	input.SpecieHaverstTime = specieHaverstTime.HarvestTime
	plant, err := h.CreatePlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "success", "Planta criada com sucesso", plant)
}

func (h *PlantHandler) DeletePlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.DeletePlantUseCaseInputDTO
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
	input.UserID = userId
	err = h.DeletePlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Planta deletada com sucesso", nil)
}

func (h *PlantHandler) FindAllPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.FindAllPlantUseCaseInputDTO
	auth := r.Header.Get("Authorization")
	userId, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	input.UserId = userId
	plants, err := h.FindAllPlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Plantas encontradas", plants)
}

func (h *PlantHandler) FindByCategoryNamePlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.FindByCategoryNamePlantUseCaseInputDTO
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
	plants, err := h.FindByCategoryNamePlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Plantas encontradas", plants)
}

func (h *PlantHandler) FindByIdPlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.FindByIdPlantUseCaseInputDTO
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
	input.UserID = userId
	plant, err := h.FindByIdPlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Planta encontrada", plant)
}

func (h *PlantHandler) FindByNamePlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.FindByNamePlantUseCaseInputDTO
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
	plant, err := h.FindByNamePlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Planta encontrada", plant)
}

func (h *PlantHandler) FindBySpecieNamePlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.FindBySpecieNamePlantUseCaseInputDTO
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
	plants, err := h.FindBySpecieNamePlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Plantas encontradas", plants)
}

func (h *PlantHandler) UpdatePlantHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_plant.UpdatePlantUseCaseInputDTO
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
	input.UserID = userId
	plant, err := h.UpdatePlantUseCase.Execute(r.Context(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Planta atualizada com sucesso", plant)
}
