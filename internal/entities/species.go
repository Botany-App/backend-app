package entities

import (
	"context"
	"time"
)

type Specie struct {
	ID                  string    `json:"id" db:"id"`
	CommonName          string    `json:"common_name" db:"common_name"`
	SpecieDescription   string    `json:"specie_description" db:"specie_description"`
	ScientificName      string    `json:"scientific_name" db:"scientific_name"`
	BotanicalFamily     string    `json:"botanical_family" db:"botanical_family"`
	GrowthType          string    `json:"growth_type" db:"growth_type"`
	IdealTemperature    float64   `json:"ideal_temperature" db:"ideal_temperature"`
	IdealClimate        string    `json:"ideal_climate" db:"ideal_climate"`
	LifeCycle           string    `json:"life_cycle" db:"life_cycle"`
	PlantingSeason      string    `json:"planting_season" db:"planting_season"`
	HarvestTime         int       `json:"harvest_time" db:"harvest_time"`
	AverageHeight       float64   `json:"average_height" db:"average_height"`
	AverageWidth        float64   `json:"average_width" db:"average_width"`
	IrrigationWeight    float64   `json:"irrigation_weight" db:"irrigation_weight"`
	FertilizationWeight float64   `json:"fertilization_weight" db:"fertilization_weight"`
	SunWeight           float64   `json:"sun_weight" db:"sun_weight"`
	ImageURL            string    `json:"image_url" db:"image_url"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type SpecieRepository interface {
	FindAll(ctx context.Context) ([]*Specie, error)
	FindById(ctx context.Context, id string) (*Specie, error)
	FindByName(ctx context.Context, commonName string) ([]*Specie, error)
}
