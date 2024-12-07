package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type SpeciesRepositoryImpl struct {
	DB *sql.DB
	RD *redis.Client
}

func NewSpeciesRepository(db *sql.DB, rd *redis.Client) *SpeciesRepositoryImpl {
	return &SpeciesRepositoryImpl{
		RD: rd,
		DB: db,
	}
}

func (r *SpeciesRepositoryImpl) FindAllPG(ctx context.Context) ([]*entities.Specie, error) {
	query := `SELECT * FROM species`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var species []*entities.Specie
	for rows.Next() {
		var specie entities.Specie
		err := rows.Scan(&specie.ID, &specie.CommonName,
			&specie.SpecieDescription, &specie.ScientificName, &specie.BotanicalFamily, &specie.GrowthType,
			&specie.IdealTemperature, &specie.IdealClimate, &specie.LifeCycle,
			&specie.PlantingSeason, &specie.HarvestTime, &specie.AverageHeight,
			&specie.AverageWidth, &specie.IrrigationWeight, &specie.FertilizationWeight,
			&specie.SunWeight, &specie.ImageURL, &specie.CreatedAt, &specie.UpdatedAt)

		if err != nil {
			return nil, err
		}
		species = append(species, &specie)
	}
	return species, nil
}

func (r *SpeciesRepositoryImpl) FindAllRD(ctx context.Context) ([]*entities.Specie, error) {
	key := "species:" + "all"
	species, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return []*entities.Specie{}, nil
	} else if err != nil {
		return nil, err
	}

	var speciesList []*entities.Specie
	for _, specie := range species {
		var specieParsed entities.Specie
		err := json.Unmarshal([]byte(specie), &specieParsed)
		if err != nil {
			return nil, err
		}
		speciesList = append(speciesList, &specieParsed)
	}
	return speciesList, nil
}

func (r *SpeciesRepositoryImpl) SetAllRD(ctx context.Context, species []*entities.Specie) error {
	key := "species:" + "all"
	data := make(map[string]interface{})
	for _, specie := range species {
		value, err := json.Marshal(specie)
		if err != nil {
			return err
		}
		data[specie.ID] = value
	}

	err := r.RD.HMSet(ctx, key, data).Err()
	if err != nil {
		return err
	}

	return r.RD.Expire(ctx, key, 20*time.Minute).Err()
}

func (r *SpeciesRepositoryImpl) FindAll(ctx context.Context) ([]*entities.Specie, error) {
	species, err := r.FindAllRD(ctx)
	if err != nil {
		return nil, err
	}
	if len(species) == 0 {
		species, err = r.FindAllPG(ctx)
		if err != nil {
			return nil, err
		}
		err = r.SetAllRD(ctx, species)
		if err != nil {
			return nil, err
		}
	}
	return species, nil
}

func (r *SpeciesRepositoryImpl) FindByIDPG(ctx context.Context, id string) (*entities.Specie, error) {
	query := `SELECT * FROM species WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, id)
	var specie entities.Specie
	err := row.Scan(&specie.ID, &specie.CommonName,
		&specie.SpecieDescription, &specie.ScientificName, &specie.BotanicalFamily, &specie.GrowthType,
		&specie.IdealTemperature, &specie.IdealClimate, &specie.LifeCycle,
		&specie.PlantingSeason, &specie.HarvestTime, &specie.AverageHeight,
		&specie.AverageWidth, &specie.IrrigationWeight, &specie.FertilizationWeight,
		&specie.SunWeight, &specie.ImageURL, &specie.CreatedAt, &specie.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Nenhum registro encontrado
		}
		return nil, err
	}
	return &specie, nil
}

func (r *SpeciesRepositoryImpl) FindByIDRD(ctx context.Context, id string) (*entities.Specie, error) {
	key := "species:" + id
	specie, err := r.RD.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var specieParsed entities.Specie
	err = json.Unmarshal([]byte(specie), &specieParsed)
	if err != nil {
		return nil, err
	}
	return &specieParsed, nil
}

func (r *SpeciesRepositoryImpl) SetByIDRD(ctx context.Context, specie *entities.Specie) error {
	key := "species:" + specie.ID
	value, err := json.Marshal(specie)
	if err != nil {
		return err
	}

	err = r.RD.Set(ctx, key, value, 20*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *SpeciesRepositoryImpl) FindById(ctx context.Context, id string) (*entities.Specie, error) {
	specie, err := r.FindByIDRD(ctx, id)
	if err != nil {
		return nil, err
	}
	if specie == nil {
		specie, err = r.FindByIDPG(ctx, id)
		if err != nil {
			return nil, err
		}
		if specie == nil {
			return nil, nil
		}
		err = r.SetByIDRD(ctx, specie)
		if err != nil {
			return nil, err
		}
	}
	return specie, nil
}

func (r *SpeciesRepositoryImpl) FindByNamePG(ctx context.Context, common_name string) ([]*entities.Specie, error) {
	query := `SELECT * FROM species WHERE common_name ILIKE $1 `
	rows, err := r.DB.QueryContext(ctx, query, "%"+common_name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var species []*entities.Specie
	for rows.Next() {
		var specie entities.Specie
		err := rows.Scan(&specie.ID, &specie.CommonName,
			&specie.SpecieDescription, &specie.ScientificName, &specie.BotanicalFamily, &specie.GrowthType,
			&specie.IdealTemperature, &specie.IdealClimate, &specie.LifeCycle,
			&specie.PlantingSeason, &specie.HarvestTime, &specie.AverageHeight,
			&specie.AverageWidth, &specie.IrrigationWeight, &specie.FertilizationWeight,
			&specie.SunWeight, &specie.ImageURL, &specie.CreatedAt, &specie.UpdatedAt)
		if err != nil {
			return nil, err
		}
		species = append(species, &specie)
	}
	return species, nil
}

func (r *SpeciesRepositoryImpl) FindByNameRD(ctx context.Context, name string) ([]*entities.Specie, error) {
	key := "species:" + "name:" + name
	species, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return []*entities.Specie{}, nil
	} else if err != nil {
		return nil, err
	}

	var speciesList []*entities.Specie
	for _, specie := range species {
		var specieParsed entities.Specie
		err := json.Unmarshal([]byte(specie), &specieParsed)
		if err != nil {
			return nil, err
		}
		speciesList = append(speciesList, &specieParsed)
	}
	return speciesList, nil
}

func (r *SpeciesRepositoryImpl) SetByNameRD(ctx context.Context, name string, species []*entities.Specie) error {
	for _, specie := range species {
		key := "species:" + "name:" + name

		value, err := json.Marshal(specie)
		if err != nil {
			return err
		}

		err = r.RD.Set(ctx, key, value, 20*time.Minute).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *SpeciesRepositoryImpl) FindByName(ctx context.Context, common_name string) ([]*entities.Specie, error) {
	species, err := r.FindByNameRD(ctx, common_name)
	if err != nil {
		return nil, err
	}
	if len(species) == 0 {
		species, err = r.FindByNamePG(ctx, common_name)
		if err != nil {
			return nil, err
		}
		err = r.SetByNameRD(ctx, common_name, species)
		if err != nil {
			return nil, err
		}
	}
	return species, nil
}
