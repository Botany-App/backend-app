package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type PlantRepositoryImpl struct {
	DB *sql.DB
	RD *redis.Client
}

func NewPlantRepositoryImpl(db *sql.DB, rd *redis.Client) *PlantRepositoryImpl {
	return &PlantRepositoryImpl{
		DB: db,
		RD: rd,
	}
}

func (r *PlantRepositoryImpl) Create(ctx context.Context, plant *entities.Plant) (string, error) {
	query := `INSERT INTO plants (id, plant_name, plant_description, planting_date,
	estimated_harvest_date, plant_status, curreting_height, currenting_width,
	 irrigation_week, health_status, last_irrigation,  last_fertilization,
	 sun_exposure, fertilization_week, user_id UUID, species_id) 
	 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`

	idParsed, err := uuid.Parse(plant.ID)
	if err != nil {
		return "", err
	}
	userIdParsed, err := uuid.Parse(plant.UserID)
	if err != nil {
		return "", err
	}

	specieId, err := uuid.Parse(plant.SpeciesID)
	if err != nil {
		return "", err
	}

	_, err = r.DB.ExecContext(ctx, query, idParsed, plant.PlantName, plant.PlantDescription,
		plant.PlantingDate, plant.EstimatedHarvestDate, plant.PlantStatus, plant.CurrentHeight,
		plant.CurrentWidth, plant.IrrigationWeek, plant.HealthStatus, plant.LastFertilization,
		plant.LastFertilization, plant.SunExposure, plant.FertilizationWeek, userIdParsed, specieId)
	if err != nil {
		return "", err
	}

	query = `INSERT INTO plant_categories (id, plant_id, category_id) VALUES ($1,$2,$3)`
	for _, categoryItem := range plant.CategoriesPlant {
		categoryId, err := uuid.Parse(categoryItem)
		if err != nil {
			return "", err
		}
		_, err = r.DB.ExecContext(ctx, query, uuid.New(), idParsed, categoryId)
		if err != nil {
			r.DB.ExecContext(ctx, `DELETE FROM plants WHERE id = $1`, idParsed)
			return "", err
		}
	}
	key := "plant:" + plant.UserID
	plantData, err := json.Marshal(plant)
	if err != nil {
		return "", err
	}

	_, err = r.RD.HSet(ctx, key, plant.ID, plantData).Result()
	if err != nil {
		return "", err
	}
	return plant.ID, nil
}

func (r *PlantRepositoryImpl) FindByIDPG(ctx context.Context, userId, id string) (*entities.Plant, error) {
	query := `SELECT * FROM plants WHERE id = $1 AND user_id = $2`

	var plant entities.Plant
	err := r.DB.QueryRowContext(ctx, query, id, userId).Scan(
		&plant.ID, &plant.PlantName, &plant.PlantDescription, &plant.PlantingDate,
		&plant.EstimatedHarvestDate, &plant.PlantStatus, &plant.CurrentHeight, &plant.CurrentWidth,
		&plant.IrrigationWeek, &plant.HealthStatus, &plant.LastIrrigation, &plant.LastFertilization,
		&plant.SunExposure, &plant.FertilizationWeek, &plant.UserID, &plant.SpeciesID,
		&plant.CreatedAt, &plant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &plant, nil
}

func (r *PlantRepositoryImpl) FindByIDRD(ctx context.Context, userId, id string) (*entities.Plant, error) {
	key := "plant:" + userId + ":" + id
	plantData, err := r.RD.HGet(ctx, key, id).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria no Redis: %w", err)
	}

	var plant entities.Plant
	err = json.Unmarshal([]byte(plantData), &plant)
	if err != nil {
		return nil, err
	}

	return &plant, nil
}

func (r *PlantRepositoryImpl) FindByID(ctx context.Context, userId, id string) (*entities.Plant, error) {
	plant, err := r.FindByIDRD(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if plant == nil {
		plant, err = r.FindByIDPG(ctx, userId, id)
		if err != nil {
			return nil, err
		}
		if plant == nil {
			return nil, nil
		}

		key := "plant:" + userId + ":" + id
		plantData, err := json.Marshal(plant)
		if err != nil {
			return nil, err
		}

		_, err = r.RD.HSet(ctx, key, id, plantData).Result()
		if err != nil {
			return nil, err
		}
	}

	return plant, nil
}

func (r *PlantRepositoryImpl) FindBySpeciesNamePG(ctx context.Context, userId, speciesName string) ([]*entities.Plant, error) {
	query := `
		SELECT p.* FROM plants p
		JOIN species s ON p.species_id = s.id
		WHERE s.common_name ILIKE %$1% AND p.user_id = $2`

	rows, err := r.DB.QueryContext(ctx, query, speciesName, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []*entities.Plant
	for rows.Next() {
		var plant entities.Plant
		err := rows.Scan(
			&plant.ID, &plant.PlantName, &plant.PlantDescription, &plant.PlantingDate,
			&plant.EstimatedHarvestDate, &plant.PlantStatus, &plant.CurrentHeight, &plant.CurrentWidth,
			&plant.IrrigationWeek, &plant.HealthStatus, &plant.LastIrrigation, &plant.LastFertilization,
			&plant.SunExposure, &plant.FertilizationWeek, &plant.UserID, &plant.SpeciesID,
			&plant.CreatedAt, &plant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindBySpeciesNameRD(ctx context.Context, userId, speciesName string) ([]*entities.Plant, error) {
	key := "plants:" + userId + ":species:" + speciesName
	plantsData, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar plantas no Redis: %w", err)
	}

	var plants []*entities.Plant
	for _, plantData := range plantsData {
		var plant entities.Plant
		err := json.Unmarshal([]byte(plantData), &plant)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindBySpeciesName(ctx context.Context, userId, speciesName string) ([]*entities.Plant, error) {
	// Primeiro tenta buscar no cache (Redis)
	plants, err := r.FindBySpeciesNameRD(ctx, userId, speciesName)
	if err != nil {
		return nil, err
	}
	if plants != nil {
		// Se já encontrou no cache, retorna os dados do Redis
		return plants, nil
	}

	// Se não encontrou no Redis, busca no banco de dados (PostgreSQL)
	plants, err = r.FindBySpeciesNamePG(ctx, userId, speciesName)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}

	// Armazena no Redis para futuras consultas
	key := "plants:" + userId + ":species:" + speciesName
	for _, plant := range plants {
		plantData, err := json.Marshal(plant)
		if err != nil {
			return nil, err
		}

		_, err = r.RD.HSet(ctx, key, plant.ID, plantData).Result()
		if err != nil {
			return nil, err
		}
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindByCategoryNamePG(ctx context.Context, userId, categoryName string) ([]*entities.Plant, error) {
	query := `
		SELECT p.* FROM plants p
		JOIN plant_categories pc ON p.id = pc.plant_id
		JOIN categories c ON pc.category_id = c.id
		WHERE c.name ILIKE %$1% AND p.user_id = $2`

	rows, err := r.DB.QueryContext(ctx, query, categoryName, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []*entities.Plant
	for rows.Next() {
		var plant entities.Plant
		err := rows.Scan(
			&plant.ID, &plant.PlantName, &plant.PlantDescription, &plant.PlantingDate,
			&plant.EstimatedHarvestDate, &plant.PlantStatus, &plant.CurrentHeight, &plant.CurrentWidth,
			&plant.IrrigationWeek, &plant.HealthStatus, &plant.LastIrrigation, &plant.LastFertilization,
			&plant.SunExposure, &plant.FertilizationWeek, &plant.UserID, &plant.SpeciesID,
			&plant.CreatedAt, &plant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindByCategoryNameRD(ctx context.Context, userId, categoryName string) ([]*entities.Plant, error) {
	key := "plants:" + userId + ":category:" + categoryName
	plantsData, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar plantas no Redis: %w", err)
	}

	var plants []*entities.Plant
	for _, plantData := range plantsData {
		var plant entities.Plant
		err := json.Unmarshal([]byte(plantData), &plant)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*entities.Plant, error) {
	// Primeiro tenta buscar no cache (Redis)
	plants, err := r.FindByCategoryNameRD(ctx, userId, categoryName)
	if err != nil {
		return nil, err
	}
	if plants != nil {
		// Se já encontrou no cache, retorna os dados do Redis
		return plants, nil
	}

	// Se não encontrou no Redis, busca no banco de dados (PostgreSQL)
	plants, err = r.FindByCategoryNamePG(ctx, userId, categoryName)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}

	// Armazena no Redis para futuras consultas
	key := "plants:" + userId + ":category:" + categoryName
	for _, plant := range plants {
		plantData, err := json.Marshal(plant)
		if err != nil {
			return nil, err
		}

		_, err = r.RD.HSet(ctx, key, plant.ID, plantData).Result()
		if err != nil {
			return nil, err
		}
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindByNamePG(ctx context.Context, userId, name string) ([]*entities.Plant, error) {
	query := `SELECT * FROM plants WHERE plant_name ILIKE $1 AND user_id = $2`

	rows, err := r.DB.QueryContext(ctx, query, "%"+name+"%", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []*entities.Plant
	for rows.Next() {
		var plant entities.Plant
		err := rows.Scan(
			&plant.ID, &plant.PlantName, &plant.PlantDescription, &plant.PlantingDate,
			&plant.EstimatedHarvestDate, &plant.PlantStatus, &plant.CurrentHeight, &plant.CurrentWidth,
			&plant.IrrigationWeek, &plant.HealthStatus, &plant.LastIrrigation, &plant.LastFertilization,
			&plant.SunExposure, &plant.FertilizationWeek, &plant.UserID, &plant.SpeciesID,
			&plant.CreatedAt, &plant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindByNameRD(ctx context.Context, userId, name string) ([]*entities.Plant, error) {
	key := "plants:" + userId + ":name:" + name
	plantsData, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar plantas no Redis: %w", err)
	}

	var plants []*entities.Plant
	for _, plantData := range plantsData {
		var plant entities.Plant
		err := json.Unmarshal([]byte(plantData), &plant)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}
func (r *PlantRepositoryImpl) FindByName(ctx context.Context, userId, name string) ([]*entities.Plant, error) {
	// Primeiro tenta buscar no cache (Redis)
	plants, err := r.FindByNameRD(ctx, userId, name)
	if err != nil {
		return nil, err
	}
	if plants != nil {
		// Se já encontrou no cache, retorna os dados do Redis
		return plants, nil
	}

	// Se não encontrou no Redis, busca no banco de dados (PostgreSQL)
	plants, err = r.FindByNamePG(ctx, userId, name)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}

	// Armazena no Redis para futuras consultas
	key := "plants:" + userId + ":name:" + name
	for _, plant := range plants {
		plantData, err := json.Marshal(plant)
		if err != nil {
			return nil, err
		}

		_, err = r.RD.HSet(ctx, key, plant.ID, plantData).Result()
		if err != nil {
			return nil, err
		}
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindAllPG(ctx context.Context, userId string) ([]*entities.Plant, error) {
	query := `SELECT * FROM plants WHERE user_id = $1`

	rows, err := r.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []*entities.Plant
	for rows.Next() {
		var plant entities.Plant
		err := rows.Scan(
			&plant.ID, &plant.PlantName, &plant.PlantDescription, &plant.PlantingDate,
			&plant.EstimatedHarvestDate, &plant.PlantStatus, &plant.CurrentHeight, &plant.CurrentWidth,
			&plant.IrrigationWeek, &plant.HealthStatus, &plant.LastIrrigation, &plant.LastFertilization,
			&plant.SunExposure, &plant.FertilizationWeek, &plant.UserID, &plant.SpeciesID,
			&plant.CreatedAt, &plant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindAllRD(ctx context.Context, userId string) ([]*entities.Plant, error) {
	key := "plants:" + userId
	plantsData, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar plantas no Redis: %w", err)
	}

	var plants []*entities.Plant
	for _, plantData := range plantsData {
		var plant entities.Plant
		err := json.Unmarshal([]byte(plantData), &plant)
		if err != nil {
			return nil, err
		}
		plants = append(plants, &plant)
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindAll(ctx context.Context, userId string) ([]*entities.Plant, error) {
	// Tenta buscar no cache (Redis)
	plants, err := r.FindAllRD(ctx, userId)
	if err != nil {
		return nil, err
	}
	if plants != nil {
		// Se já encontrou no cache, retorna os dados do Redis
		return plants, nil
	}

	// Se não encontrou no Redis, busca no banco de dados (PostgreSQL)
	plants, err = r.FindAllPG(ctx, userId)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}

	// Armazena no Redis para futuras consultas
	key := "plants:" + userId
	for _, plant := range plants {
		plantData, err := json.Marshal(plant)
		if err != nil {
			return nil, err
		}

		_, err = r.RD.HSet(ctx, key, plant.ID, plantData).Result()
		if err != nil {
			return nil, err
		}
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) UpdatePlantPG(ctx context.Context, plant *entities.Plant) error {
	query := `
		UPDATE plants
		SET 
			plant_name = $1,
			plant_description = $2,
			planting_date = $3,
			estimated_harvest_date = $4,
			plant_status = $5,
			current_height = $6,
			current_width = $7,
			irrigation_week = $8,
			health_status = $9,
			last_irrigation = $10,
			last_fertilization = $11,
			sun_exposure = $12,
			fertilization_week = $13,
			species_id = $14,
		WHERE id = $15 AND user_id = $16
	`
	_, err := r.DB.ExecContext(ctx, query,
		plant.PlantName,
		plant.PlantDescription,
		plant.PlantingDate,
		plant.EstimatedHarvestDate,
		plant.PlantStatus,
		plant.CurrentHeight,
		plant.CurrentWidth,
		plant.IrrigationWeek,
		plant.HealthStatus,
		plant.LastIrrigation,
		plant.LastFertilization,
		plant.SunExposure,
		plant.FertilizationWeek,
		plant.SpeciesID,
		plant.ID,
		plant.UserID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar planta no PostgreSQL: %w", err)
	}

	return nil
}

func (r *PlantRepositoryImpl) UpdatePlantCache(ctx context.Context, plant *entities.Plant) error {
	keyById := "plant:" + plant.UserID + ":" + plant.ID
	keyByName := "plants:" + plant.UserID + ":name:" + plant.PlantName
	keyByCategory := "plants:" + plant.UserID + ":category:" + plant.CategoriesPlant[0]
	keyBySpecies := "plants:" + plant.UserID + ":species:" + plant.SpeciesID
	keyAllPlants := "plants:" + plant.UserID

	// Remove os caches relacionados
	if err := r.RD.Del(ctx, keyById, keyByName, keyByCategory, keyBySpecies).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache de planta por ID, Nome, Categoria ou Espécie: %w", err)
	}
	if err := r.RD.Del(ctx, keyAllPlants).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache geral de plantas: %w", err)
	}

	return nil
}

func (r *PlantRepositoryImpl) Update(ctx context.Context, plant *entities.Plant) error {
	// Atualiza no PostgreSQL
	if err := r.UpdatePlantPG(ctx, plant); err != nil {
		return err
	}

	// Atualiza/remova o cache correspondente
	if err := r.UpdatePlantCache(ctx, plant); err != nil {
		log.Printf("erro ao atualizar cache da planta %s: %v", plant.ID, err)
	}

	return nil
}

func (r *PlantRepositoryImpl) DeletePlantPG(ctx context.Context, userId, id string) error {
	query := `
		DELETE FROM plants
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, id, userId)
	if err != nil {
		return fmt.Errorf("erro ao deletar planta no PostgreSQL: %w", err)
	}

	return nil
}

func (r *PlantRepositoryImpl) DeletePlantCache(ctx context.Context, userId, id, name, categoryName, speciesName string) error {
	keyByID := "plant:" + userId + ":" + id
	keyByName := "plants:" + userId + ":name:" + name
	keyByCategory := "plants:" + userId + ":category:" + categoryName
	keyBySpecies := "plants:" + userId + ":species:" + speciesName
	keyAllPlants := "plants:" + userId

	// Remove os caches relacionados
	if err := r.RD.Del(ctx, keyByID, keyByName, keyByCategory, keyBySpecies).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache por ID, Nome, Categoria ou Espécie: %w", err)
	}
	if err := r.RD.Del(ctx, keyAllPlants).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache geral de plantas: %w", err)
	}

	return nil
}

func (r *PlantRepositoryImpl) Delete(ctx context.Context, userId, id string) error {
	// Verifica se a planta existe
	plant, err := r.FindByID(ctx, userId, id)
	if err != nil {
		return err
	}
	if plant == nil {
		return errors.New("planta não encontrada")
	}

	// Deleta no PostgreSQL
	if err := r.DeletePlantPG(ctx, userId, id); err != nil {
		return err
	}

	// Remove os caches relacionados
	if err := r.DeletePlantCache(ctx, userId, id, plant.PlantName, plant.CategoriesPlant[0], plant.SpeciesID); err != nil {
		log.Printf("erro ao remover cache da planta %s: %v", id, err)
	}

	return nil
}
