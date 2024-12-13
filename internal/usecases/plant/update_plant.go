package usecases_plant

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
)

type UpdatePlantUseCaseInputDTO struct {
	ID                   string    `json:"id"`
	PlantName            string    `json:"plant_name"`
	PlantDescription     string    `json:"plant_description"`
	PlantingDate         time.Time `json:"planting_date"`
	EstimatedHarvestDate time.Time `json:"estimated_harvest_date"`
	PlantStatus          string    `json:"plant_status"`
	CurrentHeight        float64   `json:"current_height"`
	CurrentWidth         float64   `json:"current_width"`
	IrrigationWeek       int       `json:"irrigation_week"`
	HealthStatus         string    `json:"health_status"`
	LastIrrigation       time.Time `json:"last_irrigation"`
	LastFertilization    time.Time `json:"last_fertilization"`
	SunExposure          float64   `json:"sun_exposure"`
	FertilizationWeek    float64   `json:"fertilization_week"`
	UserID               string    `json:"user_id"`
	SpeciesID            string    `json:"species_id"`
	CategoriesPlant      []string  `json:"categories_plant"`
}

type UpdatePlantUseCase struct {
	PlantRepo entities.PlantRepository
}

func NewUpdatePlantUseCase(plantRepository entities.PlantRepository) *UpdatePlantUseCase {
	return &UpdatePlantUseCase{
		PlantRepo: plantRepository,
	}
}

func (u *UpdatePlantUseCase) Execute(ctx context.Context, input UpdatePlantUseCaseInputDTO) (*entities.PlantWithCategory, error) {
	// Verificar se a planta existe
	log.Println(input.UserID, input.ID)
	existingPlant, err := u.PlantRepo.FindByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar planta: %w", err)
	}
	if existingPlant == nil {
		return nil, errors.New("planta não encontrada")
	}

	// Verificar se o nome da planta já existe (excluindo a própria planta)
	plantWithSameName, err := u.PlantRepo.FindByName(ctx, input.UserID, input.PlantName)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar nome da planta: %w", err)
	}
	for _, plant := range plantWithSameName {
		if plant.PlantId != input.ID {
			return nil, errors.New("nome da planta já está em uso")
		}
	}

	// Identificar campos alterados
	updatedFields := make(map[string]interface{})
	if existingPlant.PlantName != input.PlantName {
		updatedFields["PlantName"] = input.PlantName
	}
	if existingPlant.PlantDescription != input.PlantDescription {
		updatedFields["PlantDescription"] = input.PlantDescription
	}
	if !existingPlant.PlantingDate.Equal(input.PlantingDate) {
		updatedFields["PlantingDate"] = input.PlantingDate
	}
	if !existingPlant.EstimatedHarvestDate.Equal(input.EstimatedHarvestDate) {
		updatedFields["EstimatedHarvestDate"] = input.EstimatedHarvestDate
	}
	if existingPlant.PlantStatus != input.PlantStatus {
		updatedFields["PlantStatus"] = input.PlantStatus
	}
	if existingPlant.CurrentHeight != input.CurrentHeight {
		updatedFields["CurrentHeight"] = input.CurrentHeight
	}
	if existingPlant.CurrentWidth != input.CurrentWidth {
		updatedFields["CurrentWidth"] = input.CurrentWidth
	}
	if existingPlant.IrrigationWeek != input.IrrigationWeek {
		updatedFields["IrrigationWeek"] = input.IrrigationWeek
	}
	if existingPlant.HealthStatus != input.HealthStatus {
		updatedFields["HealthStatus"] = input.HealthStatus
	}
	if !existingPlant.LastIrrigation.Equal(input.LastIrrigation) {
		updatedFields["LastIrrigation"] = input.LastIrrigation
	}
	if !existingPlant.LastFertilization.Equal(input.LastFertilization) {
		updatedFields["LastFertilization"] = input.LastFertilization
	}
	if existingPlant.SunExposure != input.SunExposure {
		updatedFields["SunExposure"] = input.SunExposure
	}
	if existingPlant.FertilizationWeek != input.FertilizationWeek {
		updatedFields["FertilizationWeek"] = input.FertilizationWeek
	}
	if existingPlant.UserId != input.UserID {
		updatedFields["UserID"] = input.UserID
	}
	if existingPlant.SpeciesId != input.SpeciesID {
		updatedFields["SpeciesID"] = input.SpeciesID
	}
	existingCategories := extractCategoryIDs(existingPlant.Category)
	if len(existingCategories) != len(input.CategoriesPlant) || !equalStringSlices(existingCategories, input.CategoriesPlant) {
		updatedFields["CategoriesPlant"] = input.CategoriesPlant
	}

	// Atualizar planta
	updatedPlant := &entities.Plant{
		Id:                   input.ID,
		PlantName:            input.PlantName,
		PlantDescription:     input.PlantDescription,
		PlantingDate:         input.PlantingDate,
		EstimatedHarvestDate: input.EstimatedHarvestDate,
		PlantStatus:          input.PlantStatus,
		CurrentHeight:        input.CurrentHeight,
		CurrentWidth:         input.CurrentWidth,
		IrrigationWeek:       input.IrrigationWeek,
		HealthStatus:         input.HealthStatus,
		LastIrrigation:       input.LastIrrigation,
		LastFertilization:    input.LastFertilization,
		SunExposure:          input.SunExposure,
		FertilizationWeek:    input.FertilizationWeek,
		UserId:               input.UserID,
		SpeciesId:            input.SpeciesID,
		CategoriesPlant:      input.CategoriesPlant,
		UpdatedAt:            time.Now(),
	}

	if err := u.PlantRepo.Update(ctx, updatedPlant); err != nil {
		return nil, fmt.Errorf("erro ao atualizar planta: %w", err)
	}

	newPlant, err := u.PlantRepo.FindByID(ctx, input.UserID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar planta atualizada: %w", err)
	}

	return newPlant, nil
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)
	for _, v := range a {
		m[v]++
	}
	for _, v := range b {
		if m[v] == 0 {
			return false
		}
		m[v]--
	}
	return true
}

func extractCategoryIDs(categories []entities.CategoryPlant) []string {
	ids := make([]string, len(categories))
	for i, category := range categories {
		ids[i] = category.Id // Substitua por como o ID é representado na struct CategoryPlant
	}
	return ids
}
