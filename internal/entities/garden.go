package entities

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Garden struct {
	Id                 string    `json:"id"`
	UserId             string    `json:"user_id"`
	GardenName         string    `json:"garden_name"`
	GardenDescription  string    `json:"garden_description"`
	GardenLocation     string    `json:"garden_location"`
	TotalArea          float64   `json:"total_area"`
	CurrentingHeigth   float64   `json:"currenting_heigth"`
	CurrentingWidth    float64   `json:"currenting_width"`
	PlantingDate       time.Time `json:"planting_date"`
	LastIrrigation     time.Time `json:"last_irrigation"`
	LastFertilization  time.Time `json:"last_fertilization"`
	Irrigation_week    int       `json:"irrigation_week"`
	Sun_Exposure       int       `json:"sun_exposure"`
	Fertilization_week int       `json:"fertilization_week"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CategoriesPlantId  []string  `json:"categories_plant"`
	PlantsId           []string  `json:"plants_id"`
}

type GardenOutputDTO struct {
	Id                 string          `json:"id"`
	UserId             string          `json:"user_id"`
	GardenName         string          `json:"garden_name"`
	GardenDescription  string          `json:"garden_description"`
	GardenLocation     string          `json:"garden_location"`
	TotalArea          float64         `json:"total_area"`
	CurrentingHeigth   float64         `json:"currenting_heigth"`
	CurrentingWidth    float64         `json:"currenting_width"`
	PlantingDate       time.Time       `json:"planting_date"`
	LastIrrigation     time.Time       `json:"last_irrigation"`
	LastFertilization  time.Time       `json:"last_fertilization"`
	Irrigation_week    int             `json:"irrigation_week"`
	Sun_Exposure       int             `json:"sun_exposure"`
	Fertilization_week int             `json:"fertilization_week"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	CategoriesPlant    []CategoryPlant `json:"categories_plant"`
	Plants             []Plant         `json:"plants"`
}

type GardenRepository interface {
	Create(ctx context.Context, garden *Garden) (string, error)
	FindByID(ctx context.Context, userId, id string) (*GardenOutputDTO, error)
	FindByName(ctx context.Context, userId, name string) ([]*GardenOutputDTO, error)
	FindByLocation(ctx context.Context, userId, location string) ([]*GardenOutputDTO, error)
	FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*GardenOutputDTO, error)
	FindAll(ctx context.Context, userId string) ([]*GardenOutputDTO, error)
}

func NewGarden(
	name, location, description string,
	area, heigth, width float64,
	plantingDate, lastIrrigation, lastFertilization time.Time,
	irrigationWeek, sunExposure, fertilizationWeek int,
	categoriesPlantId, plantsId []string,
) (*Garden, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if location == "" {
		return nil, errors.New("location is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	if area == 0 {
		return nil, errors.New("area is required")
	}
	if heigth == 0 {
		return nil, errors.New("heigth is required")
	}
	if width == 0 {
		return nil, errors.New("width is required")
	}
	if plantingDate.IsZero() {
		plantingDate = time.Now()
	}
	if lastIrrigation.IsZero() {
		lastIrrigation = time.Now()
	}
	if lastFertilization.IsZero() {
		lastFertilization = time.Now()
	}
	if irrigationWeek == 0 {
		return nil, errors.New("irrigation week is required")
	}
	if sunExposure == 0 {
		return nil, errors.New("sun exposure is required")
	}
	if fertilizationWeek == 0 {
		return nil, errors.New("fertilization week is required")
	}

	return &Garden{
		Id:                 uuid.New().String(),
		GardenName:         name,
		GardenLocation:     location,
		GardenDescription:  description,
		TotalArea:          area,
		CurrentingHeigth:   heigth,
		CurrentingWidth:    width,
		PlantingDate:       plantingDate,
		LastIrrigation:     lastIrrigation,
		LastFertilization:  lastFertilization,
		Irrigation_week:    irrigationWeek,
		Sun_Exposure:       sunExposure,
		Fertilization_week: fertilizationWeek,
		CategoriesPlantId:  categoriesPlantId,
		PlantsId:           plantsId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}
