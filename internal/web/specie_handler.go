package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases_specie "github.com/lucasBiazon/botany-back/internal/usecases/specie"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type SpecieHandler struct {
	FindAllSpeciesUseCase   *usecases_specie.FindAllSpecieUseCase
	FindByIdSpecieUseCase   *usecases_specie.FindByIdSpecieUseCase
	FindByNameSpecieUseCase *usecases_specie.FindByNameSpecieUseCase
}

func NewSpecieHandler(findAllSpeciesUseCase *usecases_specie.FindAllSpecieUseCase, findByIdSpecieUseCase *usecases_specie.FindByIdSpecieUseCase, findByNameSpecieUseCase *usecases_specie.FindByNameSpecieUseCase) *SpecieHandler {
	return &SpecieHandler{
		FindAllSpeciesUseCase:   findAllSpeciesUseCase,
		FindByIdSpecieUseCase:   findByIdSpecieUseCase,
		FindByNameSpecieUseCase: findByNameSpecieUseCase,
	}
}

func (h *SpecieHandler) FindAllSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	_, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	species, err := h.FindAllSpeciesUseCase.Execute(context.Background())
	if len(species) == 0 {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Nenhuma espécie encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar espécies", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Espécies encontradas", species)
}

func (h *SpecieHandler) FindByIdSpecieHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	_, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var input usecases_specie.FindByIdSpecieInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	species, err := h.FindByIdSpecieUseCase.Execute(context.Background(), input)
	if species == nil {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Espécie não encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar espécie", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Espécie encontrada", species)
}

func (h *SpecieHandler) FindByNameSpecieHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	_, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.JsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var input usecases_specie.FindByNameSpecieInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	species, err := h.FindByNameSpecieUseCase.Execute(context.Background(), input)
	if len(species) == 0 {
		utils.JsonResponse(w, http.StatusNotFound, "error", "Espécie não encontrada", nil)
		return
	}
	if err != nil {
		utils.JsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar espécie", err.Error())
		return
	}
	utils.JsonResponse(w, http.StatusOK, "success", "Espécie encontrada", species)
}
