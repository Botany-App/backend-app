package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_garden "github.com/lucasBiazon/botany-back/internal/usecases/garden"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type GardenHandler struct {
	CreateGardenUseCase             *usecases_garden.CreateGardenUseCase
	DeleteGardenUseCase             *usecases_garden.DeleteGardenUseCase
	FindAllGardenUseCase            *usecases_garden.FindAllGardenUseCase
	FindByIdGardenUseCase           *usecases_garden.FindByIdGardenUseCase
	UpdateGardenUseCase             *usecases_garden.UpdateGardenUseCase
	FindByLocationGardenUseCase     *usecases_garden.FindByLocationGardenUseCase
	FindByNameGardenUseCase         *usecases_garden.FindByNameGardenUseCase
	FindByCategoryNameGardenUseCase *usecases_garden.FindByCategoryNameGardenUseCase
}

func NewGardenHandler(
	createGardenUseCase *usecases_garden.CreateGardenUseCase,
	deleteGardenUseCase *usecases_garden.DeleteGardenUseCase,
	findAllGardenUseCase *usecases_garden.FindAllGardenUseCase,
	findByIdGardenUseCase *usecases_garden.FindByIdGardenUseCase,
	updateGardenUseCase *usecases_garden.UpdateGardenUseCase,
	findByLocationGardenUseCase *usecases_garden.FindByLocationGardenUseCase,
	findByNameGardenUseCase *usecases_garden.FindByNameGardenUseCase,
	findByCategoryNameGardenUseCase *usecases_garden.FindByCategoryNameGardenUseCase,
) *GardenHandler {
	return &GardenHandler{
		CreateGardenUseCase:             createGardenUseCase,
		DeleteGardenUseCase:             deleteGardenUseCase,
		FindAllGardenUseCase:            findAllGardenUseCase,
		FindByIdGardenUseCase:           findByIdGardenUseCase,
		UpdateGardenUseCase:             updateGardenUseCase,
		FindByLocationGardenUseCase:     findByLocationGardenUseCase,
		FindByNameGardenUseCase:         findByNameGardenUseCase,
		FindByCategoryNameGardenUseCase: findByCategoryNameGardenUseCase,
	}
}

func (h *GardenHandler) CreateGardenHandler(w http.ResponseWriter, r *http.Request) {

	var input usecases_garden.CreateGardenUseCaseInputDTO
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
	garden, err := h.CreateGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusCreated, "success", "Jardim criado com sucesso", garden)
}

func (h *GardenHandler) DeleteGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.DeleteGardenUseCaseInputDTO
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
	err = h.DeleteGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardim deletado com sucesso", nil)
}

func (h *GardenHandler) FindAllGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.FindAllGardenUseCaseInputDTO
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
	gardens, err := h.FindAllGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardins encontrados", gardens)
}

func (h *GardenHandler) FindByIdGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.FindByIdGardenUseCaseInputDTO
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
	garden, err := h.FindByIdGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardim encontrado", garden)
}

func (h *GardenHandler) UpdateGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.UpdateGardenUseCaseInputDTO
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
	garden, err := h.UpdateGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardim atualizado com sucesso", garden)
}

func (h *GardenHandler) FindByLocationGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.FindByLocationGardenUseCaseInputDTO
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
	gardens, err := h.FindByLocationGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardins encontrados", gardens)
}

func (h *GardenHandler) FindByNameGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.FindByNameGardenUseCaseInputDTO
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
	gardens, err := h.FindByNameGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardins encontrados", gardens)
}

func (h *GardenHandler) FindByCategoryNameGardenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases_garden.FindByCategoryNameGardenUseCaseInputDTO
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
	gardens, err := h.FindByCategoryNameGardenUseCase.Execute(context.Background(), input)
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Jardins encontrados", gardens)
}
