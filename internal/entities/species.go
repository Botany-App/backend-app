package entities

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Species struct {
	ID                  uuid.UUID `json:"id"`
	NameSpecies         string    `json:"name_species"`
	DescriptionSpecies  string    `json:"description_species"`
	FertilizationWeight float64   `json:"fertilization_weight"`
	SunWeight           float64   `json:"sun_weight"`
	IrrigationWeight    float64   `json:"irrigation_weight"`
	TimeToHarvest       float64   `json:"time_to_harvest"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type SpeciesRepository interface {
	GetAll(ctx context.Context) ([]*Species, error)
	GetByName(ctx context.Context, name string) ([]*Species, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Species, error)
	GetByHarvestTimeRange(ctx context.Context, min, max float64) ([]*Species, error)
	GetByWeights(ctx context.Context, minSun, maxSun, minFert, maxFert float64) ([]*Species, error)
}
