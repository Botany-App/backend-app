package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

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
	idParse, err := uuid.Parse(plant.Id)
	if err != nil {
		return "", fmt.Errorf("erro ao converter ID da planta: %w", err)
	}

	userIdParse, err := uuid.Parse(plant.UserId)
	if err != nil {
		return "", fmt.Errorf("erro ao converter ID do usuário: %w", err)
	}

	specieIdParse, err := uuid.Parse(plant.SpeciesId)
	if err != nil {
		return "", fmt.Errorf("erro ao converter ID da espécie: %w", err)
	}

	query := `INSERT INTO plants (id, plant_name, plant_description, planting_date,
		estimated_harvest_date, plant_status, current_height, current_width,
		irrigation_week, health_status, last_irrigation, last_fertilization,
		sun_exposure, fertilization_week, user_id, species_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err = r.DB.ExecContext(ctx, query, idParse, plant.PlantName, plant.PlantDescription,
		plant.PlantingDate, plant.EstimatedHarvestDate, plant.PlantStatus, plant.CurrentHeight,
		plant.CurrentWidth, plant.IrrigationWeek, plant.HealthStatus, plant.LastIrrigation,
		plant.LastFertilization, plant.SunExposure, plant.FertilizationWeek, userIdParse, specieIdParse)
	if err != nil {
		return "", fmt.Errorf("erro ao inserir planta: %w", err)
	}

	if len(plant.CategoriesPlant) != 0 {
		query = `INSERT INTO plant_categories (id, plant_id, category_id) VALUES ($1, $2, $3)`
		for _, categoryItem := range plant.CategoriesPlant {
			categoryId, err := uuid.Parse(categoryItem)
			if err != nil {
				return "", fmt.Errorf("erro ao converter ID da categoria: %w", err)
			}
			_, err = r.DB.ExecContext(ctx, query, uuid.New(), plant.Id, categoryId)
			if err != nil {
				return "", fmt.Errorf("erro ao inserir categoria: %w", err)
			}
		}
	}
	return plant.Id, nil
}

func (r *PlantRepositoryImpl) FindByIDPG(ctx context.Context, userId, id string) (*entities.PlantWithCategory, error) {

	idParse, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	query := `SELECT 
		p.id AS plant_id,
		p.plant_name,
		p.plant_description,
		p.planting_date,
		p.estimated_harvest_date,
		p.plant_status,
		p.current_height,
		p.current_width,
		p.irrigation_week,
		p.health_status,
		p.last_irrigation,
		p.last_fertilization,
		p.sun_exposure,
		p.fertilization_week,
		p.user_id,
		p.species_id,
		p.created_at AS plant_created_at,
		p.updated_at AS plant_updated_at,
		cp.id AS category_id,
		cp.category_name
	FROM 
		plants p
	LEFT JOIN 
		plant_categories pc ON p.id = pc.plant_id
	LEFT JOIN 
		categories_plants cp ON pc.category_id = cp.id
	WHERE 
		p.user_id = $1 AND p.id = $2;`

	rows, err := r.DB.Query(query, userIdParse, idParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	var plant *entities.PlantWithCategory
	var categories []entities.CategoryPlant

	for rows.Next() {
		var category entities.CategoryPlant
		tempPlant := entities.PlantWithCategory{}

		err := rows.Scan(
			&tempPlant.PlantId,
			&tempPlant.PlantName,
			&tempPlant.PlantDescription,
			&tempPlant.PlantingDate,
			&tempPlant.EstimatedHarvestDate,
			&tempPlant.PlantStatus,
			&tempPlant.CurrentHeight,
			&tempPlant.CurrentWidth,
			&tempPlant.IrrigationWeek,
			&tempPlant.HealthStatus,
			&tempPlant.LastIrrigation,
			&tempPlant.LastFertilization,
			&tempPlant.SunExposure,
			&tempPlant.FertilizationWeek,
			&tempPlant.UserId,
			&tempPlant.SpeciesId,
			&tempPlant.PlantCreatedAt,
			&tempPlant.PlantUpdatedAt,
			&category.Id,
			&category.Name,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		if plant == nil {
			plant = &tempPlant
		}

		if category.Id != "" {
			categories = append(categories, category)
		}
	}
	if plant != nil {
		plant.Category = categories
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if plant == nil {
		log.Println("Nenhuma planta encontrada para os parâmetros fornecidos.")
	}
	return plant, nil
}

func (r *PlantRepositoryImpl) FindByID(ctx context.Context, userId, id string) (*entities.PlantWithCategory, error) {
	plant, err := r.FindByIDPG(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	if plant == nil {
		return nil, nil
	}
	return plant, nil
}

func (r *PlantRepositoryImpl) FindBySpeciesNamePG(ctx context.Context, userId, speciesName string) ([]*entities.PlantWithCategory, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT 
			p.id AS plant_id,
			p.plant_name,
			p.plant_description,
			p.planting_date,
			p.estimated_harvest_date,
			p.plant_status,
			p.current_height,
			p.current_width,
			p.irrigation_week,
			p.health_status,
			p.last_irrigation,
			p.last_fertilization,
			p.sun_exposure,
			p.fertilization_week,
			p.user_id,
			p.species_id,
			p.created_at AS plant_created_at,
			p.updated_at AS plant_updated_at,
			cp.id AS category_id,
			cp.category_name
	FROM 
			plants p
	JOIN 
			species s ON p.species_id = s.id
	LEFT JOIN 
			plant_categories pc ON p.id = pc.plant_id
	LEFT JOIN 
			categories_plants cp ON pc.category_id = cp.id
	WHERE 
			s.common_name ILIKE $1 AND p.user_id = $2;
	`

	rows, err := r.DB.Query(query, speciesName, userIdParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	plantMap := make(map[string]*entities.PlantWithCategory)

	for rows.Next() {
		var category entities.CategoryPlant
		tempPlant := entities.PlantWithCategory{}

		err := rows.Scan(
			&tempPlant.PlantId,
			&tempPlant.PlantName,
			&tempPlant.PlantDescription,
			&tempPlant.PlantingDate,
			&tempPlant.EstimatedHarvestDate,
			&tempPlant.PlantStatus,
			&tempPlant.CurrentHeight,
			&tempPlant.CurrentWidth,
			&tempPlant.IrrigationWeek,
			&tempPlant.HealthStatus,
			&tempPlant.LastIrrigation,
			&tempPlant.LastFertilization,
			&tempPlant.SunExposure,
			&tempPlant.FertilizationWeek,
			&tempPlant.UserId,
			&tempPlant.SpeciesId,
			&tempPlant.PlantCreatedAt,
			&tempPlant.PlantUpdatedAt,
			&category.Id,
			&category.Name,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingPlant, exists := plantMap[tempPlant.PlantId]
		if !exists {
			existingPlant = &tempPlant
			existingPlant.Category = []entities.CategoryPlant{}
			plantMap[tempPlant.PlantId] = existingPlant
		}

		if category.Id != "" {
			existingPlant.Category = append(existingPlant.Category, category)
		}
	}

	var plantsWithCategories []*entities.PlantWithCategory
	for _, plant := range plantMap {
		plantsWithCategories = append(plantsWithCategories, plant)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(plantsWithCategories) == 0 {
		log.Println("Nenhuma planta encontrada para os parâmetros fornecidos.")
	}
	return plantsWithCategories, nil
}

func (r *PlantRepositoryImpl) FindBySpeciesName(ctx context.Context, userId, speciesName string) ([]*entities.PlantWithCategory, error) {
	plants, err := r.FindBySpeciesNamePG(ctx, userId, speciesName)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindByCategoryNamePG(ctx context.Context, userId, categoryName string) ([]*entities.PlantWithCategory, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT 
    p.id AS plant_id,
    p.plant_name,
    p.plant_description,
    p.planting_date,
    p.estimated_harvest_date,
    p.plant_status,
    p.current_height,
    p.current_width,
    p.irrigation_week,
    p.health_status,
    p.last_irrigation,
    p.last_fertilization,
    p.sun_exposure,
    p.fertilization_week,
    p.user_id,
    p.species_id,
    p.created_at AS plant_created_at,
    p.updated_at AS plant_updated_at,
    c.id AS category_id,
    c.category_name
FROM 
    plants p
JOIN 
    plant_categories pc ON p.id = pc.plant_id
JOIN 
    categories_plants c ON pc.category_id = c.id
WHERE 
    p.user_id = $1 AND c.category_name ILIKE $2;
	`

	rows, err := r.DB.Query(query, userIdParse, "%"+categoryName+"%")
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	plantMap := make(map[string]*entities.PlantWithCategory)

	for rows.Next() {
		var category entities.CategoryPlant
		tempPlant := entities.PlantWithCategory{}

		err := rows.Scan(
			&tempPlant.PlantId,
			&tempPlant.PlantName,
			&tempPlant.PlantDescription,
			&tempPlant.PlantingDate,
			&tempPlant.EstimatedHarvestDate,
			&tempPlant.PlantStatus,
			&tempPlant.CurrentHeight,
			&tempPlant.CurrentWidth,
			&tempPlant.IrrigationWeek,
			&tempPlant.HealthStatus,
			&tempPlant.LastIrrigation,
			&tempPlant.LastFertilization,
			&tempPlant.SunExposure,
			&tempPlant.FertilizationWeek,
			&tempPlant.UserId,
			&tempPlant.SpeciesId,
			&tempPlant.PlantCreatedAt,
			&tempPlant.PlantUpdatedAt,
			&category.Id,
			&category.Name,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingPlant, exists := plantMap[tempPlant.PlantId]
		if !exists {
			existingPlant = &tempPlant
			existingPlant.Category = []entities.CategoryPlant{}
			plantMap[tempPlant.PlantId] = existingPlant
		}

		existingPlant.Category = append(existingPlant.Category, category)
	}

	var plantsWithCategories []*entities.PlantWithCategory
	for _, plant := range plantMap {
		plantsWithCategories = append(plantsWithCategories, plant)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(plantsWithCategories) == 0 {
		log.Println("Nenhuma planta encontrada para os parâmetros fornecidos.")
	}

	return plantsWithCategories, nil
}

func (r *PlantRepositoryImpl) FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*entities.PlantWithCategory, error) {

	plants, err := r.FindByCategoryNamePG(ctx, userId, categoryName)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}
	return plants, nil
}

func (r *PlantRepositoryImpl) FindByNamePG(ctx context.Context, userId, plantName string) ([]*entities.PlantWithCategory, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	query := `
		SELECT 
			p.id AS plant_id,
			p.plant_name,
			p.plant_description,
			p.planting_date,
			p.estimated_harvest_date,
			p.plant_status,
			p.current_height,
			p.current_width,
			p.irrigation_week,
			p.health_status,
			p.last_irrigation,
			p.last_fertilization,
			p.sun_exposure,
			p.fertilization_week,
			p.user_id,
			p.species_id,
			p.created_at AS plant_created_at,
			p.updated_at AS plant_updated_at,
			c.id AS category_id,
			c.category_name
	FROM 
			plants p
	LEFT JOIN 
			plant_categories pc ON p.id = pc.plant_id
	LEFT JOIN 
			categories_plants c ON pc.category_id = c.id
	WHERE 
			p.plant_name ILIKE $1 AND p.user_id = $2;
	`

	rows, err := r.DB.Query(query, "%"+plantName+"%", userIdParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	plantMap := make(map[string]*entities.PlantWithCategory)

	for rows.Next() {
		var category entities.CategoryPlant
		tempPlant := entities.PlantWithCategory{}

		err := rows.Scan(
			&tempPlant.PlantId,
			&tempPlant.PlantName,
			&tempPlant.PlantDescription,
			&tempPlant.PlantingDate,
			&tempPlant.EstimatedHarvestDate,
			&tempPlant.PlantStatus,
			&tempPlant.CurrentHeight,
			&tempPlant.CurrentWidth,
			&tempPlant.IrrigationWeek,
			&tempPlant.HealthStatus,
			&tempPlant.LastIrrigation,
			&tempPlant.LastFertilization,
			&tempPlant.SunExposure,
			&tempPlant.FertilizationWeek,
			&tempPlant.UserId,
			&tempPlant.SpeciesId,
			&tempPlant.PlantCreatedAt,
			&tempPlant.PlantUpdatedAt,
			&category.Id,
			&category.Name,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingPlant, exists := plantMap[tempPlant.PlantId]
		if !exists {
			existingPlant = &tempPlant
			existingPlant.Category = []entities.CategoryPlant{}
			plantMap[tempPlant.PlantId] = existingPlant
		}

		if category.Id != "" {
			existingPlant.Category = append(existingPlant.Category, category)
		}
	}

	var plantsWithCategories []*entities.PlantWithCategory
	for _, plant := range plantMap {
		plantsWithCategories = append(plantsWithCategories, plant)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(plantsWithCategories) == 0 {
		log.Println("Nenhuma planta encontrada para os parâmetros fornecidos.")
	}

	return plantsWithCategories, nil
}

func (r *PlantRepositoryImpl) FindByName(ctx context.Context, userId, name string) ([]*entities.PlantWithCategory, error) {
	plants, err := r.FindByNamePG(ctx, userId, name)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}

	return plants, nil
}

func (r *PlantRepositoryImpl) FindAllPG(ctx context.Context, userId string) ([]*entities.PlantWithCategory, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	query := `
	SELECT 
			p.id AS plant_id,
			p.plant_name,
			p.plant_description,
			p.planting_date,
			p.estimated_harvest_date,
			p.plant_status,
			p.current_height,
			p.current_width,
			p.irrigation_week,
			p.health_status,
			p.last_irrigation,
			p.last_fertilization,
			p.sun_exposure,
			p.fertilization_week,
			p.user_id,
			p.species_id,
			p.created_at AS plant_created_at,
			p.updated_at AS plant_updated_at,
			c.id AS category_id,
			c.category_name
	FROM 
			plants p
	LEFT JOIN 
			plant_categories pc ON p.id = pc.plant_id
	LEFT JOIN 
			categories_plants c ON pc.category_id = c.id
	WHERE 
			p.user_id = $1;`

	rows, err := r.DB.Query(query, userIdParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	plantMap := make(map[string]*entities.PlantWithCategory)

	for rows.Next() {
		var category entities.CategoryPlant
		tempPlant := entities.PlantWithCategory{}

		err := rows.Scan(
			&tempPlant.PlantId,
			&tempPlant.PlantName,
			&tempPlant.PlantDescription,
			&tempPlant.PlantingDate,
			&tempPlant.EstimatedHarvestDate,
			&tempPlant.PlantStatus,
			&tempPlant.CurrentHeight,
			&tempPlant.CurrentWidth,
			&tempPlant.IrrigationWeek,
			&tempPlant.HealthStatus,
			&tempPlant.LastIrrigation,
			&tempPlant.LastFertilization,
			&tempPlant.SunExposure,
			&tempPlant.FertilizationWeek,
			&tempPlant.UserId,
			&tempPlant.SpeciesId,
			&tempPlant.PlantCreatedAt,
			&tempPlant.PlantUpdatedAt,
			&category.Id,
			&category.Name,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingPlant, exists := plantMap[tempPlant.PlantId]
		if !exists {
			existingPlant = &tempPlant
			existingPlant.Category = []entities.CategoryPlant{}
			plantMap[tempPlant.PlantId] = existingPlant
		}

		if category.Id != "" {
			existingPlant.Category = append(existingPlant.Category, category)
		}
	}

	var plantsWithCategories []*entities.PlantWithCategory
	for _, plant := range plantMap {
		plantsWithCategories = append(plantsWithCategories, plant)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(plantsWithCategories) == 0 {
		log.Println("Nenhuma planta encontrada para o usuário fornecido.")
	}

	return plantsWithCategories, nil
}

func (r *PlantRepositoryImpl) FindAll(ctx context.Context, userId string) ([]*entities.PlantWithCategory, error) {

	plants, err := r.FindAllPG(ctx, userId)
	if err != nil {
		return nil, err
	}
	if plants == nil {
		return nil, nil
	}
	return plants, nil
}

func (r *PlantRepositoryImpl) UpdatePlantPG(ctx context.Context, plant *entities.Plant) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}

	// Atualizar campos da planta
	updatePlantQuery := `
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
            updated_at = $14
        WHERE id = $15;
    `

	_, err = tx.Exec(
		updatePlantQuery,
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
		time.Now(),
		plant.Id,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar planta: %v", err)
	}

	deleteCategoriesQuery := `
        DELETE FROM plant_categories
        WHERE plant_id = $1;
    `
	idParsed, err := uuid.Parse(plant.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao converter ID da planta: %v", err)
	}

	_, err = tx.Exec(deleteCategoriesQuery, idParsed)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao remover categorias da planta: %v", err)
	}

	insertCategoryQuery := `
        INSERT INTO plant_categories (id, plant_id, category_id)
        VALUES ($1, $2, $3);
    `

	for _, categoryId := range plant.CategoriesPlant {
		_, err = tx.Exec(insertCategoryQuery, uuid.New(), plant.Id, categoryId)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao adicionar categoria %v à planta: %v", categoryId, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

func (r *PlantRepositoryImpl) Update(ctx context.Context, plant *entities.Plant) error {
	log.Println("Iniciando processo de atualização.")
	if err := r.UpdatePlantPG(ctx, plant); err != nil {
		return err
	}
	return nil
}

func (r *PlantRepositoryImpl) DeletePlantPG(ctx context.Context, userID, plantID string) error {

	if userID == "" || plantID == "" {
		log.Println("Erro: ID do usuário ou da planta não fornecido.")
		return fmt.Errorf("ID do usuário ou da planta não fornecido")
	}
	query := `
		DELETE FROM plants
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, plantID, userID)
	if err != nil {
		log.Printf("Erro ao deletar planta no PostgreSQL: %v\n", err)
		return fmt.Errorf("erro ao deletar planta no PostgreSQL: %w", err)
	}

	log.Println("Planta deletada com sucesso no PostgreSQL.")
	return nil
}

func (r *PlantRepositoryImpl) Delete(ctx context.Context, userID, plantID string) error {

	plantWithCategory, err := r.FindByID(ctx, userID, plantID)
	if err != nil {
		log.Printf("Erro ao buscar planta %s: %v", plantID, err)
		return err
	}
	if plantWithCategory == nil {
		log.Printf("Planta %s não encontrada.", plantID)
		return errors.New("planta não encontrada")
	}

	if err := r.DeletePlantPG(ctx, userID, plantID); err != nil {
		return err
	}

	log.Println("Processo de exclusão da planta concluído com sucesso.")
	return nil
}
