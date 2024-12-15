package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type GardenRepositoryImpl struct {
	DB *sql.DB
}

func NewGardenRepository(db *sql.DB) *GardenRepositoryImpl {
	return &GardenRepositoryImpl{DB: db}
}

func (r *GardenRepositoryImpl) Create(ctx context.Context, garden *entities.Garden) (string, error) {
	idParsed, err := uuid.Parse(garden.Id)
	if err != nil {
		return "", err
	}
	userIdParsed, err := uuid.Parse(garden.UserId)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO gardens (id, user_id, garden_name, garden_description,  garden_location
	total_area, currenting_height, currenting_width, planting_date, last_irrigation, last_fertilization, 
	irrigation_week, sun_exposure, fertilization_week)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = r.DB.ExecContext(ctx, query, idParsed, userIdParsed, garden.GardenName, garden.GardenDescription, garden.GardenLocation,
		garden.TotalArea, garden.CurrentingHeight, garden.CurrentingWidth, garden.PlantingDate, garden.LastIrrigation, garden.LastFertilization,
		garden.IrrigationWeek, garden.SunExposure, garden.FertilizationWeek)
	if err != nil {
		return "", err
	}

	if len(garden.CategoriesPlantId) > 0 {
		query = `INSERT INTO garden_categories (id, garden_id, category_id) VALUES ($1, $2, $3)`
		for _, categoryPlantId := range garden.CategoriesPlantId {
			categoryId, err := uuid.Parse(categoryPlantId)
			if err != nil {
				return "", err
			}
			_, err = r.DB.ExecContext(ctx, query, uuid.New(), idParsed, categoryId)
			if err != nil {
				return "", err
			}
		}
	}

	if len(garden.PlantsId) > 0 {
		query = `INSERT INTO garden_plant (id, garden_id, plant_id) VALUES ($1, $2, $3)`
		for _, plantId := range garden.PlantsId {
			plantIdParsed, err := uuid.Parse(plantId)
			if err != nil {
				return "", err
			}
			_, err = r.DB.ExecContext(ctx, query, uuid.New(), idParsed, plantIdParsed)
			if err != nil {
				return "", err
			}
		}
	}

	return garden.Id, nil
}

func (r *GardenRepositoryImpl) FindByID(ctx context.Context, userId, id string) (*entities.GardenOutputDTO, error) {
	idParse, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	query := `SELECT 
		g.id AS garden_id,
		g.garden_name,
		g.garden_description,
		g.garden_location,
		g.total_area,
		g.currenting_height,
		g.currenting_width,
		g.planting_date,
		g.last_irrigation,
		g.last_fertilization,
		g.irrigation_week,
		g.sun_exposure,
		g.fertilization_week,
		g.created_at AS garden_created_at,
		g.updated_at AS garden_updated_at,
		gc.category_id,
		cp.category_name,
		gp.plant_id,
		p.plant_name
	FROM 
		gardens g
	LEFT JOIN 
		garden_categories gc ON g.id = gc.garden_id
	LEFT JOIN 
		categories_plants cp ON gc.category_id = cp.id
	LEFT JOIN 
		garden_plant gp ON g.id = gp.garden_id
	LEFT JOIN 
		plants p ON gp.plant_id = p.id
	WHERE 
		g.user_id = $1 AND g.id = $2;`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, idParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	var garden *entities.GardenOutputDTO
	var categories []entities.CategoryPlant
	var plants []entities.Plant

	for rows.Next() {
		var category entities.CategoryPlant
		var plant entities.Plant
		tempGarden := entities.GardenOutputDTO{}

		err := rows.Scan(
			&tempGarden.Id,
			&tempGarden.GardenName,
			&tempGarden.GardenDescription,
			&tempGarden.GardenLocation,
			&tempGarden.TotalArea,
			&tempGarden.CurrentingHeight,
			&tempGarden.CurrentingWidth,
			&tempGarden.PlantingDate,
			&tempGarden.LastIrrigation,
			&tempGarden.LastFertilization,
			&tempGarden.IrrigationWeek,
			&tempGarden.SunExposure,
			&tempGarden.FertilizationWeek,
			&tempGarden.CreatedAt,
			&tempGarden.UpdatedAt,
			&category.Id,
			&category.Name,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		if garden == nil {
			garden = &tempGarden
		}

		if category.Id != "" {
			categories = append(categories, category)
		}
		if plant.Id != "" {
			plants = append(plants, plant)
		}
	}

	if garden != nil {
		garden.CategoriesPlant = categories
		garden.Plants = plants
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if garden == nil {
		log.Println("Nenhum jardim encontrado para os parâmetros fornecidos.")
	}
	return garden, nil
}

func (r *GardenRepositoryImpl) FindAll(ctx context.Context, userId string) ([]*entities.GardenOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	query := `SELECT 
		g.id AS garden_id,
		g.garden_name,
		g.garden_description,
		g.garden_location,
		g.total_area,
		g.currenting_height,
		g.currenting_width,
		g.planting_date,
		g.last_irrigation,
		g.last_fertilization,
		g.irrigation_week,
		g.sun_exposure,
		g.fertilization_week,
		g.created_at AS garden_created_at,
		g.updated_at AS garden_updated_at,
		gc.category_id,
		cp.category_name,
		gp.plant_id,
		p.plant_name
	FROM 
		gardens g
	LEFT JOIN 
		garden_categories gc ON g.id = gc.garden_id
	LEFT JOIN 
		categories_plants cp ON gc.category_id = cp.id
	LEFT JOIN 
		garden_plant gp ON g.id = gp.garden_id
	LEFT JOIN 
		plants p ON gp.plant_id = p.id
	WHERE 
		 g.user_id = $1;`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	gardenMap := make(map[string]*entities.GardenOutputDTO)

	for rows.Next() {
		var category entities.CategoryPlant
		var plant entities.Plant
		tempGarden := entities.GardenOutputDTO{}

		err := rows.Scan(
			&tempGarden.Id,
			&tempGarden.GardenName,
			&tempGarden.GardenDescription,
			&tempGarden.GardenLocation,
			&tempGarden.TotalArea,
			&tempGarden.CurrentingHeight,
			&tempGarden.CurrentingWidth,
			&tempGarden.PlantingDate,
			&tempGarden.LastIrrigation,
			&tempGarden.LastFertilization,
			&tempGarden.IrrigationWeek,
			&tempGarden.SunExposure,
			&tempGarden.FertilizationWeek,
			&tempGarden.CreatedAt,
			&tempGarden.UpdatedAt,
			&category.Id,
			&category.Name,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingGarden, exists := gardenMap[tempGarden.Id]
		if !exists {
			existingGarden = &tempGarden
			existingGarden.CategoriesPlant = []entities.CategoryPlant{}
			existingGarden.Plants = []entities.Plant{}
			gardenMap[tempGarden.Id] = existingGarden
		}

		if category.Id != "" {
			existingGarden.CategoriesPlant = append(existingGarden.CategoriesPlant, category)
		}
		if plant.Id != "" {
			existingGarden.Plants = append(existingGarden.Plants, plant)
		}
	}

	var gardens []*entities.GardenOutputDTO
	for _, garden := range gardenMap {
		gardens = append(gardens, garden)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(gardens) == 0 {
		log.Println("Nenhum jardim encontrado para os parâmetros fornecidos.")
	}

	return gardens, nil
}

func (r *GardenRepositoryImpl) FindByLocation(ctx context.Context, userId, gardenLocation string) ([]*entities.GardenOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	query := `SELECT 
		g.id AS garden_id,
		g.garden_name,
		g.garden_description,
		g.garden_location,
		g.total_area,
		g.currenting_height,
		g.currenting_width,
		g.planting_date,
		g.last_irrigation,
		g.last_fertilization,
		g.irrigation_week,
		g.sun_exposure,
		g.fertilization_week,
		g.created_at AS garden_created_at,
		g.updated_at AS garden_updated_at,
		gc.category_id,
		cp.category_name,
		gp.plant_id,
		p.plant_name
	FROM 
		gardens g
	LEFT JOIN 
		garden_categories gc ON g.id = gc.garden_id
	LEFT JOIN 
		categories_plants cp ON gc.category_id = cp.id
	LEFT JOIN 
		garden_plant gp ON g.id = gp.garden_id
	LEFT JOIN 
		plants p ON gp.plant_id = p.id
	WHERE 
		g.garden_location ILIKE $1 AND g.user_id = $2;`

	rows, err := r.DB.QueryContext(ctx, query, "%"+gardenLocation+"%", userIdParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	gardenMap := make(map[string]*entities.GardenOutputDTO)

	for rows.Next() {
		var category entities.CategoryPlant
		var plant entities.Plant
		tempGarden := entities.GardenOutputDTO{}

		err := rows.Scan(
			&tempGarden.Id,
			&tempGarden.GardenName,
			&tempGarden.GardenDescription,
			&tempGarden.GardenLocation,
			&tempGarden.TotalArea,
			&tempGarden.CurrentingHeight,
			&tempGarden.CurrentingWidth,
			&tempGarden.PlantingDate,
			&tempGarden.LastIrrigation,
			&tempGarden.LastFertilization,
			&tempGarden.IrrigationWeek,
			&tempGarden.SunExposure,
			&tempGarden.FertilizationWeek,
			&tempGarden.CreatedAt,
			&tempGarden.UpdatedAt,
			&category.Id,
			&category.Name,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingGarden, exists := gardenMap[tempGarden.Id]
		if !exists {
			existingGarden = &tempGarden
			existingGarden.CategoriesPlant = []entities.CategoryPlant{}
			existingGarden.Plants = []entities.Plant{}
			gardenMap[tempGarden.Id] = existingGarden
		}

		if category.Id != "" {
			existingGarden.CategoriesPlant = append(existingGarden.CategoriesPlant, category)
		}
		if plant.Id != "" {
			existingGarden.Plants = append(existingGarden.Plants, plant)
		}
	}

	var gardens []*entities.GardenOutputDTO
	for _, garden := range gardenMap {
		gardens = append(gardens, garden)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(gardens) == 0 {
		log.Println("Nenhum jardim encontrado para os parâmetros fornecidos.")
	}

	return gardens, nil
}

func (r *GardenRepositoryImpl) FindByName(ctx context.Context, userId, gardenName string) ([]*entities.GardenOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	query := `SELECT 
		g.id AS garden_id,
		g.garden_name,
		g.garden_description,
		g.garden_location,
		g.total_area,
		g.currenting_height,
		g.currenting_width,
		g.planting_date,
		g.last_irrigation,
		g.last_fertilization,
		g.irrigation_week,
		g.sun_exposure,
		g.fertilization_week,
		g.created_at AS garden_created_at,
		g.updated_at AS garden_updated_at,
		gc.category_id,
		cp.category_name,
		gp.plant_id,
		p.plant_name
	FROM 
		gardens g
	LEFT JOIN 
		garden_categories gc ON g.id = gc.garden_id
	LEFT JOIN 
		categories_plants cp ON gc.category_id = cp.id
	LEFT JOIN 
		garden_plant gp ON g.id = gp.garden_id
	LEFT JOIN 
		plants p ON gp.plant_id = p.id
	WHERE 
		g.garden_name ILIKE $1 AND g.user_id = $2;`

	rows, err := r.DB.QueryContext(ctx, query, "%"+gardenName+"%", userIdParse)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	gardenMap := make(map[string]*entities.GardenOutputDTO)

	for rows.Next() {
		var category entities.CategoryPlant
		var plant entities.Plant
		tempGarden := entities.GardenOutputDTO{}

		err := rows.Scan(
			&tempGarden.Id,
			&tempGarden.GardenName,
			&tempGarden.GardenDescription,
			&tempGarden.GardenLocation,
			&tempGarden.TotalArea,
			&tempGarden.CurrentingHeight,
			&tempGarden.CurrentingWidth,
			&tempGarden.PlantingDate,
			&tempGarden.LastIrrigation,
			&tempGarden.LastFertilization,
			&tempGarden.IrrigationWeek,
			&tempGarden.SunExposure,
			&tempGarden.FertilizationWeek,
			&tempGarden.CreatedAt,
			&tempGarden.UpdatedAt,
			&category.Id,
			&category.Name,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingGarden, exists := gardenMap[tempGarden.Id]
		if !exists {
			existingGarden = &tempGarden
			existingGarden.CategoriesPlant = []entities.CategoryPlant{}
			existingGarden.Plants = []entities.Plant{}
			gardenMap[tempGarden.Id] = existingGarden
		}

		if category.Id != "" {
			existingGarden.CategoriesPlant = append(existingGarden.CategoriesPlant, category)
		}
		if plant.Id != "" {
			existingGarden.Plants = append(existingGarden.Plants, plant)
		}
	}

	var gardens []*entities.GardenOutputDTO
	for _, garden := range gardenMap {
		gardens = append(gardens, garden)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(gardens) == 0 {
		log.Println("Nenhum jardim encontrado para os parâmetros fornecidos.")
	}

	return gardens, nil
}

func (r *GardenRepositoryImpl) FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*entities.GardenOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	query := `SELECT 
		g.id AS garden_id,
		g.garden_name,
		g.garden_description,
		g.garden_location,
		g.total_area,
		g.currenting_height,
		g.currenting_width,
		g.planting_date,
		g.last_irrigation,
		g.last_fertilization,
		g.irrigation_week,
		g.sun_exposure,
		g.fertilization_week,
		g.created_at AS garden_created_at,
		g.updated_at AS garden_updated_at,
		gc.category_id,
		cp.category_name,
		gp.plant_id,
		p.plant_name
	FROM 
		gardens g
	LEFT JOIN 
		garden_categories gc ON g.id = gc.garden_id
	LEFT JOIN 
		categories_plants cp ON gc.category_id = cp.id
	LEFT JOIN 
		garden_plant gp ON g.id = gp.garden_id
	LEFT JOIN 
		plants p ON gp.plant_id = p.id
	WHERE 
		 g.user_id = $1 AND cp.category_name ILIKE $2;`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, "%"+categoryName+"%")
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	gardenMap := make(map[string]*entities.GardenOutputDTO)

	for rows.Next() {
		var category entities.CategoryPlant
		var plant entities.Plant
		tempGarden := entities.GardenOutputDTO{}

		err := rows.Scan(
			&tempGarden.Id,
			&tempGarden.GardenName,
			&tempGarden.GardenDescription,
			&tempGarden.GardenLocation,
			&tempGarden.TotalArea,
			&tempGarden.CurrentingHeight,
			&tempGarden.CurrentingWidth,
			&tempGarden.PlantingDate,
			&tempGarden.LastIrrigation,
			&tempGarden.LastFertilization,
			&tempGarden.IrrigationWeek,
			&tempGarden.SunExposure,
			&tempGarden.FertilizationWeek,
			&tempGarden.CreatedAt,
			&tempGarden.UpdatedAt,
			&category.Id,
			&category.Name,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			log.Fatalf("Erro ao escanear resultados: %v", err)
		}

		existingGarden, exists := gardenMap[tempGarden.Id]
		if !exists {
			existingGarden = &tempGarden
			existingGarden.CategoriesPlant = []entities.CategoryPlant{}
			existingGarden.Plants = []entities.Plant{}
			gardenMap[tempGarden.Id] = existingGarden
		}

		if category.Id != "" {
			existingGarden.CategoriesPlant = append(existingGarden.CategoriesPlant, category)
		}
		if plant.Id != "" {
			existingGarden.Plants = append(existingGarden.Plants, plant)
		}
	}

	var gardens []*entities.GardenOutputDTO
	for _, garden := range gardenMap {
		gardens = append(gardens, garden)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if len(gardens) == 0 {
		log.Println("Nenhum jardim encontrado para os parâmetros fornecidos.")
	}

	return gardens, nil
}

func (r *GardenRepositoryImpl) Update(ctx context.Context, garden *entities.Garden) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	updateGardenQuery := `UPDATE gardens SET garden_name = $1, garden_description = $2, garden_location = $3,
	total_area = $4, currenting_height = $5, currenting_width = $6, planting_date = $7, last_irrigation = $8,
	last_fertilization = $9, irrigation_week = $10, sun_exposure = $11, fertilization_week = $12
	WHERE id = $13 AND user_id = $14;`

	_, err = tx.ExecContext(ctx, updateGardenQuery, garden.GardenName, garden.GardenDescription, garden.GardenLocation,
		garden.TotalArea, garden.CurrentingHeight, garden.CurrentingWidth, garden.PlantingDate, garden.LastIrrigation,
		garden.LastFertilization, garden.IrrigationWeek, garden.SunExposure, garden.FertilizationWeek, garden.Id, garden.UserId)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar jardim: %v", err)
	}

	deleteCategoriesQuery := `DELETE FROM garden_categories WHERE garden_id = $1;`
	_, err = tx.ExecContext(ctx, deleteCategoriesQuery, garden.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao deletar categorias: %v", err)
	}
	deletePlantQuery := `DELETE FROM garden_plant WHERE garden_id = $1;`
	_, err = tx.ExecContext(ctx, deletePlantQuery, garden.Id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao deletar plantas: %v", err)
	}

	if len(garden.CategoriesPlantId) > 0 {
		insertCategoriesQuery := `INSERT INTO garden_categories (id, garden_id, category_id) VALUES ($1, $2, $3);`
		for _, categoryPlantId := range garden.CategoriesPlantId {
			categoryId, err := uuid.Parse(categoryPlantId)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("erro ao converter id da categoria: %v", err)
			}
			_, err = tx.ExecContext(ctx, insertCategoriesQuery, uuid.New(), garden.Id, categoryId)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("erro ao inserir categoria: %v", err)
			}
		}
	}

	if len(garden.PlantsId) > 0 {
		insertPlantQuery := `INSERT INTO garden_plant (id, garden_id, plant_id) VALUES ($1, $2, $3);`
		for _, plantId := range garden.PlantsId {
			plantIdParsed, err := uuid.Parse(plantId)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("erro ao converter id da planta: %v", err)
			}
			_, err = tx.ExecContext(ctx, insertPlantQuery, uuid.New(), garden.Id, plantIdParsed)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("erro ao inserir planta: %v", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao commitar transação: %v", err)
	}

	return nil

}

func (r *GardenRepositoryImpl) Delete(ctx context.Context, userId, id string) error {
	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return err
	}
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM gardens WHERE id = $1 AND user_id = $2;`
	_, err = r.DB.ExecContext(ctx, query, idParsed, userIdParsed)
	if err != nil {
		return err
	}

	return nil
}
