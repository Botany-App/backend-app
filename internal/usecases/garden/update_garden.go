package usecases_garden

import (
	"context"
	"errors"
	"time"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UpdateGardenUseCaseInputDTO struct {
	ID                string    `json:"id"`
	GardenName        string    `json:"garden_name"`
	UserID            string    `json:"user_id"`
	GardenDescription string    `json:"garden_description"`
	GardenLocation    string    `json:"garden_location"`
	TotalArea         float64   `json:"total_area"`
	CurrentingHeight  float64   `json:"currenting_heigth"`
	CurrentingWidth   float64   `json:"currenting_width"`
	PlantingDate      time.Time `json:"planting_date"`
	LastIrrigation    time.Time `json:"last_irrigation"`
	LastFertilization time.Time `json:"last_fertilization"`
	IrrigationWeek    int       `json:"irrigation_week"`
	SunExposure       int       `json:"sun_exposure"`
	FertilizationWeek int       `json:"fertilization_week"`
	CategoriesPlantId []string  `json:"categories_plant"`
	PlantsId          []string  `json:"plants_id"`
}

type UpdateGardenUseCase struct {
	Repository entities.GardenRepository
}

func NewUpdateGardenUseCase(repository entities.GardenRepository) *UpdateGardenUseCase {
	return &UpdateGardenUseCase{
		Repository: repository,
	}
}

func (uc *UpdateGardenUseCase) Execute(ctx context.Context, input UpdateGardenUseCaseInputDTO) (*entities.GardenOutputDTO, error) {
	existingGarden, err := uc.Repository.FindByID(ctx, input.UserID, input.GardenName)
	if err != nil {
		return nil, err
	}
	if existingGarden == nil {
		return nil, errors.New("jardim não encontrado")
	}

	gardenName, err := uc.Repository.FindByName(ctx, input.UserID, input.GardenName)
	if err != nil {
		return nil, err
	}

	for _, garden := range gardenName {
		if garden.Id != input.ID {
			return nil, errors.New("nome do jardim já está em uso")
		}
	}

	updatedFields := make(map[string]interface{})
	if existingGarden.GardenName != input.GardenName {
		updatedFields["garden_name"] = input.GardenName
	}
	if existingGarden.GardenDescription != input.GardenDescription {
		updatedFields["garden_description"] = input.GardenDescription
	}
	if existingGarden.GardenLocation != input.GardenLocation {
		updatedFields["garden_location"] = input.GardenLocation
	}
	if existingGarden.TotalArea != input.TotalArea {
		updatedFields["total_area"] = input.TotalArea
	}
	if existingGarden.CurrentingHeight != input.CurrentingHeight {
		updatedFields["currenting_height"] = input.CurrentingHeight
	}
	if existingGarden.CurrentingWidth != input.CurrentingWidth {
		updatedFields["currenting_width"] = input.CurrentingWidth
	}
	if !existingGarden.PlantingDate.Equal(input.PlantingDate) {
		updatedFields["planting_date"] = input.PlantingDate
	}
	if !existingGarden.LastIrrigation.Equal(input.LastIrrigation) {
		updatedFields["last_irrigation"] = input.LastIrrigation
	}
	if existingGarden.LastFertilization.Equal(input.LastFertilization) {
		updatedFields["last_fertilization"] = input.LastFertilization
	}
	if existingGarden.IrrigationWeek != input.IrrigationWeek {
		updatedFields["irrigation_week"] = input.IrrigationWeek
	}
	if existingGarden.SunExposure != input.SunExposure {
		updatedFields["sun_exposure"] = input.SunExposure
	}
	if existingGarden.FertilizationWeek != input.FertilizationWeek {
		updatedFields["fertilization_week"] = input.FertilizationWeek
	}

	existingCategories := extractCategoryIDs(existingGarden.CategoriesPlant)
	if len(existingCategories) != len(input.CategoriesPlantId) || !equalStringSlices(existingCategories, input.CategoriesPlantId) {
		updatedFields["categories_plant"] = input.CategoriesPlantId
	}

	existingPlants := extractPlantIDs(existingGarden.Plants)
	if len(existingPlants) != len(input.PlantsId) || !equalStringSlices(existingPlants, input.PlantsId) {
		updatedFields["plants_id"] = input.PlantsId
	}

	updatedGarden := &entities.Garden{
		Id:                input.ID,
		GardenName:        input.GardenName,
		GardenDescription: input.GardenDescription,
		GardenLocation:    input.GardenLocation,
		TotalArea:         input.TotalArea,
		CurrentingHeight:  input.CurrentingHeight,
		CurrentingWidth:   input.CurrentingWidth,
		PlantingDate:      input.PlantingDate,
		LastIrrigation:    input.LastIrrigation,
		LastFertilization: input.LastFertilization,
		IrrigationWeek:    input.IrrigationWeek,
		SunExposure:       input.SunExposure,
		FertilizationWeek: input.FertilizationWeek,
		CategoriesPlantId: input.CategoriesPlantId,
		PlantsId:          input.PlantsId,
		UpdatedAt:         time.Now(),
		UserId:            input.UserID,
	}
	if err := uc.Repository.Update(ctx, updatedGarden); err != nil {
		return nil, err
	}

	garden, err := uc.Repository.FindByID(ctx, input.UserID, input.GardenName)
	if err != nil {
		return nil, err
	}
	return garden, nil
}

func extractCategoryIDs(categories []entities.CategoryPlant) []string {
	ids := make([]string, len(categories))
	for i, category := range categories {
		ids[i] = category.Id // Substitua por como o ID é representado na struct CategoryPlant
	}
	return ids
}

func extractPlantIDs(plants []entities.Plant) []string {
	ids := make([]string, len(plants))
	for i, plant := range plants {
		ids[i] = plant.Id // Substitua por como o ID é representado na struct Plant
	}
	return ids
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
