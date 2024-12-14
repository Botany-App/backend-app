package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CategoryTaskRepositoryImpl struct {
	DB *sql.DB
	RD *redis.Client
}

func NewCategoryTaskRepository(db *sql.DB, rd *redis.Client) *CategoryTaskRepositoryImpl {
	return &CategoryTaskRepositoryImpl{
		DB: db,
		RD: rd,
	}
}

func (r *CategoryTaskRepositoryImpl) Create(ctx context.Context, categoryTask *entities.CategoryTask) (string, error) {
	query := `INSERT INTO categories_tasks (id, category_name, category_description, user_id) VALUES ($1, $2, $3, $4)`

	idParsed, err := uuid.Parse(categoryTask.Id)
	if err != nil {
		return "", err
	}
	userIdParsed, err := uuid.Parse(categoryTask.UserId)
	if err != nil {
		return "", err
	}

	_, err = r.DB.ExecContext(ctx, query, idParsed, categoryTask.Name, categoryTask.Description, userIdParsed)
	if err != nil {
		return "", err
	}

	key := "categories_tasks:" + categoryTask.UserId
	categoryData, err := json.Marshal(categoryTask)
	if err != nil {
		return "", err
	}

	_, err = r.RD.HSet(ctx, key, categoryTask.Id, categoryData).Result()
	if err != nil {
		return "", err
	}

	return categoryTask.Id, nil
}

func (r *CategoryTaskRepositoryImpl) FindAllPG(ctx context.Context, userId string) ([]*entities.CategoryTask, error) {
	query := `SELECT id, category_name, category_description, user_id, created_at, updated_at FROM categories_tasks WHERE user_id = $1`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, query, userIdParsed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.CategoryTask
	for rows.Next() {
		var category entities.CategoryTask
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.UserId, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) FindAllFromCache(ctx context.Context, userId string) ([]*entities.CategoryTask, error) {
	key := "categories_tasks:" + userId

	categories, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return []*entities.CategoryTask{}, nil // Retorna slice vazio, sem erro
	} else if err != nil {
		return nil, err // Outros erros de Redis
	}

	var categoriesParsed []*entities.CategoryTask
	for _, category := range categories {
		var categoryParsed entities.CategoryTask
		err := json.Unmarshal([]byte(category), &categoryParsed)
		if err != nil {
			return nil, err
		}
		categoriesParsed = append(categoriesParsed, &categoryParsed)
	}

	return categoriesParsed, nil
}

func (r *CategoryTaskRepositoryImpl) CacheAllCategories(ctx context.Context, userId string, categories []*entities.CategoryTask) error {
	key := "categories_tasks:" + userId

	data := make(map[string]interface{})
	for _, category := range categories {
		categoryParsed, err := json.Marshal(category)
		if err != nil {
			return err
		}
		data[category.Id] = categoryParsed
	}

	_, err := r.RD.HSet(ctx, key, data).Result()
	if err != nil {
		return err
	}

	return r.RD.Expire(ctx, key, 20*time.Minute).Err() // Configura TTL
}

func (r *CategoryTaskRepositoryImpl) FindAll(ctx context.Context, userId string) ([]*entities.CategoryTask, error) {
	categories, err := r.FindAllFromCache(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		categories, err = r.FindAllPG(ctx, userId)
		if err != nil {
			return nil, err
		}
		if categories == nil {
			categories = []*entities.CategoryTask{}
		}

		err = r.CacheAllCategories(ctx, userId, categories)
		if err != nil {
			return nil, err
		}
	}

	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) FindByNamePG(ctx context.Context, userId, name string) ([]*entities.CategoryTask, error) {
	query := `SELECT id, category_name, category_description, user_id, created_at, updated_at 
	          FROM categories_tasks WHERE user_id = $1 AND category_name ILIKE $2`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, query, userIdParsed, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.CategoryTask
	for rows.Next() {
		var category entities.CategoryTask
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.UserId, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) FindByNameRD(ctx context.Context, userId, name string) ([]*entities.CategoryTask, error) {
	pattern := "categories_tasks:" + userId + ":" + name + "*"

	keys, err := r.RD.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var categories []*entities.CategoryTask
	for _, key := range keys {
		categoryJSON, err := r.RD.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var category entities.CategoryTask
		if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) SetCategoryByNameRD(ctx context.Context, userId string, categories []*entities.CategoryTask) error {
	for _, category := range categories {
		key := "categories_tasks:" + userId + ":" + category.Name + ":" + category.Id

		categoryJSON, err := json.Marshal(category)
		if err != nil {
			return err
		}

		err = r.RD.Set(ctx, key, categoryJSON, 10*time.Minute).Err() // Adiciona TTL
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) FindByName(ctx context.Context, userId, name string) ([]*entities.CategoryTask, error) {
	// Tenta buscar no Redis primeiro
	categories, err := r.FindByNameRD(ctx, userId, name)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 { // Cache miss, busca no PostgreSQL
		categories, err = r.FindByNamePG(ctx, userId, name)
		if err != nil {
			return nil, err
		}

		if len(categories) > 0 { // Atualiza o cache se houver resultados
			err = r.SetCategoryByNameRD(ctx, userId, categories)
			if err != nil {
				return nil, err
			}
		}
	}

	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) FindByIDPG(ctx context.Context, userId, id string) (*entities.CategoryTask, error) {
	query := `
		SELECT id, category_name, category_description, user_id, created_at, updated_at
		FROM categories_tasks
		WHERE id = $1 AND user_id = $2
	`

	// Parse dos IDs para UUID
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter userId para UUID: %w", err)
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter id para UUID: %w", err)
	}

	// Executa a consulta
	row := r.DB.QueryRowContext(ctx, query, idUUID, userUUID)

	// Mapeia os resultados para a entidade
	var category entities.CategoryTask
	if err := row.Scan(
		&category.Id,
		&category.Name,
		&category.Description,
		&category.UserId,
		&category.CreatedAt,
		&category.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Nenhum registro encontrado
		}
		return nil, fmt.Errorf("erro ao escanear categoria: %w", err)
	}

	return &category, nil
}

func (r *CategoryTaskRepositoryImpl) FindByIDRD(ctx context.Context, userId, id string) (*entities.CategoryTask, error) {
	// Define a chave para busca no Redis
	key := "categories_tasks:" + userId + ":" + id

	// Tenta recuperar os dados do Redis
	categoryJSON, err := r.RD.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria no Redis: %w", err)
	}

	// Desserializa o JSON para a entidade
	var category entities.CategoryTask
	if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
		return nil, fmt.Errorf("erro ao desserializar categoria: %w", err)
	}

	return &category, nil
}

func (r *CategoryTaskRepositoryImpl) SetCategoryByIDRD(ctx context.Context, userId, id string, category *entities.CategoryTask) error {
	// Define a chave e serializa a entidade para JSON
	key := "categories_tasks:" + userId + ":" + id

	categoryJSON, err := json.Marshal(category)
	if err != nil {
		return fmt.Errorf("erro ao serializar categoria: %w", err)
	}

	// Armazena no Redis com TTL
	err = r.RD.Set(ctx, key, categoryJSON, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("erro ao salvar categoria no Redis: %w", err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) FindById(ctx context.Context, userId, id string) (*entities.CategoryTask, error) {
	// Tenta buscar no cache do Redis
	category, err := r.FindByIDRD(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	if category == nil { // Cache miss, busca no PostgreSQL
		category, err = r.FindByIDPG(ctx, userId, id)
		if err != nil {
			return nil, err
		}

		if category != nil {
			// Atualiza o cache com a entidade encontrada
			if err := r.SetCategoryByIDRD(ctx, userId, id, category); err != nil {
				log.Printf("Erro ao atualizar cache para a categoria ID %s: %v", id, err)
			}
		}
	}

	return category, nil
}

func (r *CategoryTaskRepositoryImpl) UpdateCategoryPG(ctx context.Context, category *entities.CategoryTask) error {
	query := `
		UPDATE categories_tasks
		SET category_name = $1,
			category_description = $2
		WHERE id = $3 AND user_id = $4
	`
	_, err := r.DB.ExecContext(ctx, query,
		category.Name,
		category.Description,
		category.Id,
		category.UserId,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria no PostgreSQL: %w", err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) UpdateCategoryCache(ctx context.Context, category *entities.CategoryTask) error {
	keyById := "categories_tasks:" + category.UserId + ":" + category.Id
	keyByName := "categories_tasks:" + category.UserId + ":" + category.Name + ":" + category.Id
	keyAllCategories := "categories_tasks:" + category.UserId

	// Remove os caches relacionados
	if err := r.RD.Del(ctx, keyById, keyByName).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache por ID e Nome: %w", err)
	}
	if err := r.RD.HDel(ctx, keyAllCategories, category.Id).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache de categorias gerais: %w", err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) Update(ctx context.Context, category *entities.CategoryTask) error {
	// Atualiza no PostgreSQL
	if err := r.UpdateCategoryPG(ctx, category); err != nil {
		return err
	}

	// Atualiza/remova o cache correspondente
	if err := r.UpdateCategoryCache(ctx, category); err != nil {
		log.Printf("erro ao atualizar cache da categoria %s: %v", category.Id, err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) DeleteCategoryPG(ctx context.Context, userId, id string) error {
	query := `
		DELETE FROM categories_tasks
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, id, userId)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria no PostgreSQL: %w", err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) DeleteCategoryCache(ctx context.Context, userId, id, name string) error {
	keyByID := "categories_tasks:" + userId + ":" + id
	keyByName := "categories_tasks:" + userId + ":" + name + ":" + id
	keyAllCategories := "categories_tasks:" + userId

	// Remove os caches relacionados
	if err := r.RD.Del(ctx, keyByID, keyByName).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache por ID e Nome: %w", err)
	}
	if err := r.RD.HDel(ctx, keyAllCategories, id).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache de categorias gerais: %w", err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) Delete(ctx context.Context, userId, id string) error {
	// Verifica se a categoria existe
	category, err := r.FindByIDPG(ctx, userId, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("categoria não encontrada")
	}

	// Deleta no PostgreSQL
	if err := r.DeleteCategoryPG(ctx, userId, id); err != nil {
		return err
	}

	// Remove os caches relacionados
	if err := r.DeleteCategoryCache(ctx, userId, id, category.Name); err != nil {
		log.Printf("erro ao remover cache da categoria %s: %v", id, err)
	}

	return nil
}
