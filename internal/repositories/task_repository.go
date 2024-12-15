package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type TaskRepositoryImpl struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{DB: db}
}

func (r *TaskRepositoryImpl) Create(ctx context.Context, task *entities.Task) (string, error) {
	query := `INSERT INTO tasks (id, task_name, task_description, date_task, urgency_level, task_status, user_id)
	 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	idParsed, err := uuid.Parse(task.Id)
	if err != nil {
		return "", errors.New("invalid id")
	}

	userIdParsed, err := uuid.Parse(task.UserId)
	if err != nil {
		return "", errors.New("invalid user id")
	}

	_, err = r.DB.ExecContext(ctx, query, idParsed, task.Name, task.Description, task.TaskDate, task.UrgencyLevel, task.TaskStatus, userIdParsed)
	if err != nil {
		return "", err
	}

	if len(task.PlantsId) != 0 {
		for _, plantId := range task.PlantsId {
			_, err = r.DB.ExecContext(ctx, "INSERT INTO task_plants (id, task_id, plant_id) VALUES ($1, $2, $3)", uuid.New(), idParsed, plantId)
			if err != nil {
				return "", err
			}
		}
	}
	if len(task.GardensId) != 0 {
		for _, gardenId := range task.GardensId {
			_, err = r.DB.ExecContext(ctx, "INSERT INTO task_gardens (id, task_id, garden_id) VALUES ($1, $2, $3)", uuid.New(), idParsed, gardenId)
			if err != nil {
				return "", err
			}
		}
	}
	if len(task.CategoriesId) != 0 {
		for _, categoryId := range task.CategoriesId {
			_, err = r.DB.ExecContext(ctx, "INSERT INTO task_categories (id, task_id, category_id) VALUES ($1, $2, $3)", uuid.New(), idParsed, categoryId)
			if err != nil {
				return "", err
			}
		}
	}
	return task.Id, nil
}

func (r *TaskRepositoryImpl) Update(ctx context.Context, task *entities.Task) error {
	query := `UPDATE tasks SET task_name = $1, task_description = $2, date_task = $3, urgency_level = $4, task_status = $5 WHERE id = $6`

	idParsed, err := uuid.Parse(task.Id)
	if err != nil {
		return errors.New("invalid id")
	}

	_, err = r.DB.ExecContext(ctx, query, task.Name, task.Description, task.TaskDate, task.UrgencyLevel, task.TaskStatus, idParsed)
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, "DELETE FROM task_plants WHERE task_id = $1", idParsed)
	if err != nil {
		return err
	}
	if len(task.PlantsId) != 0 {
		for _, plantId := range task.PlantsId {
			_, err = r.DB.ExecContext(ctx, "INSERT INTO task_plants (id, task_id, plant_id) VALUES ($1, $2, $3)", uuid.New(), idParsed, plantId)
			if err != nil {
				return err
			}
		}
	}

	_, err = r.DB.ExecContext(ctx, "DELETE FROM task_gardens WHERE task_id = $1", idParsed)
	if err != nil {
		return err
	}
	if len(task.GardensId) != 0 {
		for _, gardenId := range task.GardensId {
			_, err = r.DB.ExecContext(ctx, "INSERT INTO task_gardens (id, task_id, garden_id) VALUES ($1, $2, $3)", uuid.New(), idParsed, gardenId)
			if err != nil {
				return err
			}
		}
	}

	_, err = r.DB.ExecContext(ctx, "DELETE FROM task_categories WHERE task_id = $1", idParsed)
	if err != nil {
		return err
	}
	if len(task.CategoriesId) != 0 {
		for _, categoryId := range task.CategoriesId {
			_, err = r.DB.ExecContext(ctx, "INSERT INTO task_categories (id, task_id, category_id) VALUES ($1, $2, $3)", uuid.New(), idParsed, categoryId)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *TaskRepositoryImpl) Delete(ctx context.Context, userId, id string) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`

	idParsed, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id")
	}

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("invalid user id")
	}

	_, err = r.DB.ExecContext(ctx, query, idParsed, userIdParsed)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepositoryImpl) FindByID(ctx context.Context, userId, id string) (*entities.TaskOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId: %w", err)
	}
	query := `
		SELECT 
			t.id AS task_id,
			t.task_name,
			t.task_description,
			t.date_task,
			t.urgency_level,
			t.task_status,
			t.user_id,
			t.created_at AS task_created_at,
			t.updated_at AS task_updated_at,
			tc.category_id,
			c.category_name,
			tg.garden_id,
			g.garden_name,
			tp.plant_id,
			p.plant_name
		FROM 
			tasks t
		LEFT JOIN 
			task_categories tc ON t.id = tc.task_id
		LEFT JOIN 
			categories c ON tc.category_id = c.id
		LEFT JOIN 
			task_gardens tg ON t.id = tg.task_id
		LEFT JOIN 
			gardens g ON tg.garden_id = g.id
		LEFT JOIN 
			task_plants tp ON t.id = tp.task_id
		LEFT JOIN 
			plants p ON tp.plant_id = p.id
		WHERE 
			t.user_id = $1 AND t.id = $2;
	`

	idParse, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter id: %w", err)
	}

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, idParse)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()
	var task *entities.TaskOutputDTO
	var categories []entities.CategoryTask
	var gardens []entities.Garden
	var plants []entities.Plant

	for rows.Next() {
		var category entities.CategoryTask
		var garden entities.Garden
		var plant entities.Plant
		tempTask := entities.TaskOutputDTO{}

		err := rows.Scan(
			&tempTask.Id,
			&tempTask.Name,
			&tempTask.Description,
			&tempTask.TaskDate,
			&tempTask.UrgencyLevel,
			&tempTask.TaskStatus,
			&tempTask.UserId,
			&tempTask.CreatedAt,
			&tempTask.UpdatedAt,
			&category.Id,
			&category.Name,
			&garden.Id,
			&garden.GardenName,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear resultados: %w", err)
		}

		if task == nil {
			task = &tempTask
		}

		if category.Id != "" {
			categories = append(categories, category)
		}
		if garden.Id != "" {
			gardens = append(gardens, garden)
		}
		if plant.Id != "" {
			plants = append(plants, plant)
		}
	}

	if task != nil {
		task.Categories = categories
		task.Gardens = gardens
		task.Plants = plants
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Erro durante a iteração: %v", err)
	}

	if task == nil {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (r *TaskRepositoryImpl) FindByName(ctx context.Context, userId, name string) ([]*entities.TaskOutputDTO, error) {

	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId: %w", err)
	}

	query := `
		SELECT 
			t.id AS task_id,
			t.task_name,
			t.task_description,
			t.date_task,
			t.urgency_level,
			t.task_status,
			t.user_id,
			t.created_at AS task_created_at,
			t.updated_at AS task_updated_at,
			tc.category_id,
			c.category_name,
			tg.garden_id,
			g.garden_name,
			tp.plant_id,
			p.plant_name
		FROM 
			tasks t
		LEFT JOIN 
			task_categories tc ON t.id = tc.task_id
		LEFT JOIN 
			categories c ON tc.category_id = c.id
		LEFT JOIN 
			task_gardens tg ON t.id = tg.task_id
		LEFT JOIN 
			gardens g ON tg.garden_id = g.id
		LEFT JOIN 
			task_plants tp ON t.id = tp.task_id
		LEFT JOIN 
			plants p ON tp.plant_id = p.id
		WHERE 
			t.user_id = $1 AND t.task_name ILIKE $2;
	`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()

	// Map para consolidar tarefas
	taskMap := make(map[string]*entities.TaskOutputDTO)

	for rows.Next() {
		var category entities.CategoryTask
		var garden entities.Garden
		var plant entities.Plant
		tempTask := entities.TaskOutputDTO{}

		err := rows.Scan(
			&tempTask.Id,
			&tempTask.Name,
			&tempTask.Description,
			&tempTask.TaskDate,
			&tempTask.UrgencyLevel,
			&tempTask.TaskStatus,
			&tempTask.UserId,
			&tempTask.CreatedAt,
			&tempTask.UpdatedAt,
			&category.Id,
			&category.Name,
			&garden.Id,
			&garden.GardenName,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear resultados: %w", err)
		}

		// Adiciona ou atualiza a tarefa no mapa
		existingTask, exists := taskMap[tempTask.Id]
		if !exists {
			existingTask = &tempTask
			existingTask.Categories = []entities.CategoryTask{}
			existingTask.Gardens = []entities.Garden{}
			existingTask.Plants = []entities.Plant{}
			taskMap[tempTask.Id] = existingTask
		}

		// Adiciona categorias, jardins e plantas se existirem
		if category.Id != "" {
			existingTask.Categories = append(existingTask.Categories, category)
		}
		if garden.Id != "" {
			existingTask.Gardens = append(existingTask.Gardens, garden)
		}
		if plant.Id != "" {
			existingTask.Plants = append(existingTask.Plants, plant)
		}
	}

	var tasks []*entities.TaskOutputDTO
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) FindByCategoryName(ctx context.Context, userId, categoryName string) ([]*entities.TaskOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId: %w", err)
	}

	query := `
		SELECT 
			t.id AS task_id,
			t.task_name,
			t.task_description,
			t.date_task,
			t.urgency_level,
			t.task_status,
			t.user_id,
			t.created_at AS task_created_at,
			t.updated_at AS task_updated_at,
			tc.category_id,
			c.category_name,
			tg.garden_id,
			g.garden_name,
			tp.plant_id,
			p.plant_name
		FROM 
			tasks t
		LEFT JOIN 
			task_categories tc ON t.id = tc.task_id
		LEFT JOIN 
			categories c ON tc.category_id = c.id
		LEFT JOIN 
			task_gardens tg ON t.id = tg.task_id
		LEFT JOIN 
			gardens g ON tg.garden_id = g.id
		LEFT JOIN 
			task_plants tp ON t.id = tp.task_id
		LEFT JOIN 
			plants p ON tp.plant_id = p.id
		WHERE 
			t.user_id = $1 AND c.category_name ILIKE $2;
	`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, "%"+categoryName+"%")
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()

	// Map para consolidar tarefas
	taskMap := make(map[string]*entities.TaskOutputDTO)

	for rows.Next() {
		var category entities.CategoryTask
		var garden entities.Garden
		var plant entities.Plant
		tempTask := entities.TaskOutputDTO{}

		err := rows.Scan(
			&tempTask.Id,
			&tempTask.Name,
			&tempTask.Description,
			&tempTask.TaskDate,
			&tempTask.UrgencyLevel,
			&tempTask.TaskStatus,
			&tempTask.UserId,
			&tempTask.CreatedAt,
			&tempTask.UpdatedAt,
			&category.Id,
			&category.Name,
			&garden.Id,
			&garden.GardenName,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear resultados: %w", err)
		}

		// Adiciona ou atualiza a tarefa no mapa
		existingTask, exists := taskMap[tempTask.Id]
		if !exists {
			existingTask = &tempTask
			existingTask.Categories = []entities.CategoryTask{}
			existingTask.Gardens = []entities.Garden{}
			existingTask.Plants = []entities.Plant{}
			taskMap[tempTask.Id] = existingTask
		}

		// Adiciona categorias, jardins e plantas se existirem
		if category.Id != "" {
			existingTask.Categories = append(existingTask.Categories, category)
		}
		if garden.Id != "" {
			existingTask.Gardens = append(existingTask.Gardens, garden)
		}
		if plant.Id != "" {
			existingTask.Plants = append(existingTask.Plants, plant)
		}
	}

	// Converte o mapa para uma slice
	var tasks []*entities.TaskOutputDTO
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) FindByStatus(ctx context.Context, userId, status string) ([]*entities.TaskOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId: %w", err)
	}

	query := `
		SELECT 
			t.id AS task_id,
			t.task_name,
			t.task_description,
			t.date_task,
			t.urgency_level,
			t.task_status,
			t.user_id,
			t.created_at AS task_created_at,
			t.updated_at AS task_updated_at,
			tc.category_id,
			c.category_name,
			tg.garden_id,
			g.garden_name,
			tp.plant_id,
			p.plant_name
		FROM 
			tasks t
		LEFT JOIN 
			task_categories tc ON t.id = tc.task_id
		LEFT JOIN 
			categories c ON tc.category_id = c.id
		LEFT JOIN 
			task_gardens tg ON t.id = tg.task_id
		LEFT JOIN 
			gardens g ON tg.garden_id = g.id
		LEFT JOIN 
			task_plants tp ON t.id = tp.task_id
		LEFT JOIN 
			plants p ON tp.plant_id = p.id
		WHERE 
			t.user_id = $1 AND t.task_status ILIKE $2;
	`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, "%"+status+"%")
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()

	// Map para consolidar tarefas
	taskMap := make(map[string]*entities.TaskOutputDTO)

	for rows.Next() {
		var category entities.CategoryTask
		var garden entities.Garden
		var plant entities.Plant
		tempTask := entities.TaskOutputDTO{}

		err := rows.Scan(
			&tempTask.Id,
			&tempTask.Name,
			&tempTask.Description,
			&tempTask.TaskDate,
			&tempTask.UrgencyLevel,
			&tempTask.TaskStatus,
			&tempTask.UserId,
			&tempTask.CreatedAt,
			&tempTask.UpdatedAt,
			&category.Id,
			&category.Name,
			&garden.Id,
			&garden.GardenName,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear resultados: %w", err)
		}

		// Adiciona ou atualiza a tarefa no mapa
		existingTask, exists := taskMap[tempTask.Id]
		if !exists {
			existingTask = &tempTask
			existingTask.Categories = []entities.CategoryTask{}
			existingTask.Gardens = []entities.Garden{}
			existingTask.Plants = []entities.Plant{}
			taskMap[tempTask.Id] = existingTask
		}

		// Adiciona categorias, jardins e plantas se existirem
		if category.Id != "" {
			existingTask.Categories = append(existingTask.Categories, category)
		}
		if garden.Id != "" {
			existingTask.Gardens = append(existingTask.Gardens, garden)
		}
		if plant.Id != "" {
			existingTask.Plants = append(existingTask.Plants, plant)
		}
	}

	var tasks []*entities.TaskOutputDTO
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) FindByUrgencyLevel(ctx context.Context, userId string, urgencyLevel int) ([]*entities.TaskOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId: %w", err)
	}

	query := `
		SELECT 
			t.id AS task_id,
			t.task_name,
			t.task_description,
			t.date_task,
			t.urgency_level,
			t.task_status,
			t.user_id,
			t.created_at AS task_created_at,
			t.updated_at AS task_updated_at,
			tc.category_id,
			c.category_name,
			tg.garden_id,
			g.garden_name,
			tp.plant_id,
			p.plant_name
		FROM 
			tasks t
		LEFT JOIN 
			task_categories tc ON t.id = tc.task_id
		LEFT JOIN 
			categories c ON tc.category_id = c.id
		LEFT JOIN 
			task_gardens tg ON t.id = tg.task_id
		LEFT JOIN 
			gardens g ON tg.garden_id = g.id
		LEFT JOIN 
			task_plants tp ON t.id = tp.task_id
		LEFT JOIN 
			plants p ON tp.plant_id = p.id
		WHERE 
			t.user_id = $1 AND t.urgency_level = $2;
	`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse, urgencyLevel)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()

	// Map para consolidar tarefas
	taskMap := make(map[string]*entities.TaskOutputDTO)

	for rows.Next() {
		var category entities.CategoryTask
		var garden entities.Garden
		var plant entities.Plant
		tempTask := entities.TaskOutputDTO{}

		err := rows.Scan(
			&tempTask.Id,
			&tempTask.Name,
			&tempTask.Description,
			&tempTask.TaskDate,
			&tempTask.UrgencyLevel,
			&tempTask.TaskStatus,
			&tempTask.UserId,
			&tempTask.CreatedAt,
			&tempTask.UpdatedAt,
			&category.Id,
			&category.Name,
			&garden.Id,
			&garden.GardenName,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear resultados: %w", err)
		}

		// Adiciona ou atualiza a tarefa no mapa
		existingTask, exists := taskMap[tempTask.Id]
		if !exists {
			existingTask = &tempTask
			existingTask.Categories = []entities.CategoryTask{}
			existingTask.Gardens = []entities.Garden{}
			existingTask.Plants = []entities.Plant{}
			taskMap[tempTask.Id] = existingTask
		}

		// Adiciona categorias, jardins e plantas se existirem
		if category.Id != "" {
			existingTask.Categories = append(existingTask.Categories, category)
		}
		if garden.Id != "" {
			existingTask.Gardens = append(existingTask.Gardens, garden)
		}
		if plant.Id != "" {
			existingTask.Plants = append(existingTask.Plants, plant)
		}
	}

	var tasks []*entities.TaskOutputDTO
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepositoryImpl) FindAll(ctx context.Context, userId string) ([]*entities.TaskOutputDTO, error) {
	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId: %w", err)
	}

	query := `
		SELECT 
			t.id AS task_id,
			t.task_name,
			t.task_description,
			t.date_task,
			t.urgency_level,
			t.task_status,
			t.user_id,
			t.created_at AS task_created_at,
			t.updated_at AS task_updated_at,
			tc.category_id,
			c.category_name,
			tg.garden_id,
			g.garden_name,
			tp.plant_id,
			p.plant_name
		FROM 
			tasks t
		LEFT JOIN 
			task_categories tc ON t.id = tc.task_id
		LEFT JOIN 
			categories c ON tc.category_id = c.id
		LEFT JOIN 
			task_gardens tg ON t.id = tg.task_id
		LEFT JOIN 
			gardens g ON tg.garden_id = g.id
		LEFT JOIN 
			task_plants tp ON t.id = tp.task_id
		LEFT JOIN 
			plants p ON tp.plant_id = p.id
		WHERE 
			t.user_id = $1;
	`

	rows, err := r.DB.QueryContext(ctx, query, userIdParse)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()

	// Map para consolidar tarefas
	taskMap := make(map[string]*entities.TaskOutputDTO)

	for rows.Next() {
		var category entities.CategoryTask
		var garden entities.Garden
		var plant entities.Plant
		tempTask := entities.TaskOutputDTO{}

		err := rows.Scan(
			&tempTask.Id,
			&tempTask.Name,
			&tempTask.Description,
			&tempTask.TaskDate,
			&tempTask.UrgencyLevel,
			&tempTask.TaskStatus,
			&tempTask.UserId,
			&tempTask.CreatedAt,
			&tempTask.UpdatedAt,
			&category.Id,
			&category.Name,
			&garden.Id,
			&garden.GardenName,
			&plant.Id,
			&plant.PlantName,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear resultados: %w)", err)
		}
		existingTask, exists := taskMap[tempTask.Id]
		if !exists {
			existingTask = &tempTask
			existingTask.Categories = []entities.CategoryTask{}
			existingTask.Gardens = []entities.Garden{}
			existingTask.Plants = []entities.Plant{}
			taskMap[tempTask.Id] = existingTask
		}

		// Adiciona categorias, jardins e plantas se existirem
		if category.Id != "" {
			existingTask.Categories = append(existingTask.Categories, category)
		}
		if garden.Id != "" {
			existingTask.Gardens = append(existingTask.Gardens, garden)
		}
		if plant.Id != "" {
			existingTask.Plants = append(existingTask.Plants, plant)
		}
	}

	var tasks []*entities.TaskOutputDTO
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração: %w", err)
	}

	return tasks, nil
}
