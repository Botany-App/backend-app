package entities

import "time"

type Specie struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	FertilizantionWeight float32   `json:"fertilizantion_weight"`
	Irrigation_weight    float32   `json:"irrigation_weight"`
	SumWeight            float32   `json:"sum_weight"`
	TimeToHarvest        float32   `json:"time_to_harvest"`
	CreatedAt            time.Time `json:"creation_date"`
	UpdatedAt            time.Time `json:"update_date"`
}

type SpecieRepository interface {
	GetByName(name string) (*Specie, error)
	GetByFertilizantionWeight(fertilizantion_weight float32) (*Specie, error)
	GetByIrrigationWeight(irrigation_weight float32) (*Specie, error)
	GetBySumWeight(sum_weight float32) (*Specie, error)
	GetByTimeToHarvest(time_to_harvest float32) (*Specie, error)
}
