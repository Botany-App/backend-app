package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type CategoryPlantRepositoryImpl struct {
	DB *sql.DB
	RD *redis.Client
}

func NewCategoryPlantRepository(db *sql.DB, rd *redis.Client) *CategoryPlantRepositoryImpl {
	return &CategoryPlantRepositoryImpl{
		DB: db,
		RD: rd,
	}
}

func (r *CategoryPlantRepositoryImpl) Create(ctx context.Context, categoryPlant *entities.CategoryPlant) (string, error) {
	query := `INSERT INTO categories_plants (id, category_name, category_description, user_id) VALUES ($1, $2, $3, $4)`

	idParsed, err := uuid.Parse(categoryPlant.Id)
	if err != nil {
		return "", err
	}
	userIdParsed, err := uuid.Parse(categoryPlant.UserId)
	if err != nil {
		return "", err
	}

	_, err = r.DB.ExecContext(ctx, query, idParsed, categoryPlant.Name, categoryPlant.Description, userIdParsed)
	if err != nil {
		return "", err
	}

	return categoryPlant.Id, nil
}

func (r *CategoryPlantRepositoryImpl) FindAllPG(ctx context.Context, userId string) ([]*entities.CategoryPlant, error) {
	query := `SELECT id, category_name, category_description, user_id, created_at, updated_at FROM categories_plants WHERE user_id = $1`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, query, userIdParsed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.CategoryPlant
	for rows.Next() {
		var category entities.CategoryPlant
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.UserId, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryPlantRepositoryImpl) FindAllFromCache(ctx context.Context, userId string) ([]*entities.CategoryPlant, error) {
	key := "categories_plants:" + userId

	categories, err := r.RD.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return []*entities.CategoryPlant{}, nil // Retorna slice vazio, sem erro
	} else if err != nil {
		return nil, err // Outros erros de Redis
	}

	var categoriesParsed []*entities.CategoryPlant
	for _, category := range categories {
		var categoryParsed entities.CategoryPlant
		err := json.Unmarshal([]byte(category), &categoryParsed)
		if err != nil {
			return nil, err
		}
		categoriesParsed = append(categoriesParsed, &categoryParsed)
	}

	return categoriesParsed, nil
}

func (r *CategoryPlantRepositoryImpl) CacheAllCategories(ctx context.Context, userId string, categories []*entities.CategoryPlant) error {
	key := "categories_plants:" + userId

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

	return r.RD.Expire(ctx, key, 10*time.Minute).Err() // Configura TTL
}

func (r *CategoryPlantRepositoryImpl) FindAll(ctx context.Context, userId string) ([]*entities.CategoryPlant, error) {
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
			categories = []*entities.CategoryPlant{}
		}

		err = r.CacheAllCategories(ctx, userId, categories)
		if err != nil {
			return nil, err
		}
	}

	return categories, nil
}
func (r *CategoryPlantRepositoryImpl) FindByNamePG(ctx context.Context, userId, name string) ([]*entities.CategoryPlant, error) {
	query := `SELECT id, category_name, category_description, user_id, created_at, updated_at 
	          FROM categories_plants WHERE user_id = $1 AND category_name ILIKE $2`

	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.QueryContext(ctx, query, userIdParsed, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.CategoryPlant
	for rows.Next() {
		var category entities.CategoryPlant
		err := rows.Scan(&category.Id, &category.Name, &category.Description, &category.UserId, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryPlantRepositoryImpl) FindByNameRD(ctx context.Context, userId, name string) ([]*entities.CategoryPlant, error) {
	pattern := "categories_plants:" + userId + ":" + name + "*"

	keys, err := r.RD.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var categories []*entities.CategoryPlant
	for _, key := range keys {
		categoryJSON, err := r.RD.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}

		var category entities.CategoryPlant
		if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryPlantRepositoryImpl) SetCategoryByNameRD(ctx context.Context, userId, name string, categories []*entities.CategoryPlant) error {
	for _, category := range categories {
		key := "categories_plants:" + userId + ":" + name + ":" + category.Id

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

func (r *CategoryPlantRepositoryImpl) FindByName(ctx context.Context, userId, name string) ([]*entities.CategoryPlant, error) {
	// Tenta buscar no Redis primeiro
	category, err := r.FindByNameRD(ctx, userId, name)
	if err != nil {
		return nil, err
	}

	if category == nil { // Cache miss, busca no PostgreSQL
		category, err = r.FindByNamePG(ctx, userId, name)
		if err != nil {
			return nil, err
		}

		if category != nil { // Atualiza o cache
			err = r.SetCategoryByNameRD(ctx, userId, name, category)
			if err != nil {
				return nil, err
			}
		}
	}

	return category, nil
}

func (r *CategoryPlantRepositoryImpl) FindByIDPG(ctx context.Context, userId, id string) (*entities.CategoryPlant, error) {
	query := `
			SELECT id, category_name, category_description, user_id, created_at, updated_at
			FROM categories_plants
			WHERE id = $1 AND user_id = $2
	`
	row := r.DB.QueryRowContext(ctx, query, id, userId)

	var category entities.CategoryPlant
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
		return nil, err
	}

	return &category, nil
}

func (r *CategoryPlantRepositoryImpl) FindByIDRD(ctx context.Context, userId, id string) (*entities.CategoryPlant, error) {
	key := "categories_plants:" + userId + ":" + id // Chave com id

	categoryJSON, err := r.RD.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err // Outros erros do Redis
	}

	var category entities.CategoryPlant
	if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryPlantRepositoryImpl) SetCategoryByIDRD(ctx context.Context, userId, id string, category *entities.CategoryPlant) error {
	key := "categories_plants:" + userId + ":" + id

	categoryJSON, err := json.Marshal(category)
	if err != nil {
		return err
	}

	err = r.RD.Set(ctx, key, categoryJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) FindById(ctx context.Context, userId, id string) (*entities.CategoryPlant, error) {
	// Tenta buscar no Redis primeiro
	category, err := r.FindByIDRD(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	if category == nil { // Cache miss, busca no PostgreSQL
		category, err = r.FindByIDPG(ctx, userId, id)
		if err != nil {
			return nil, err
		}

		if category != nil { // Atualiza o cache
			err = r.SetCategoryByIDRD(ctx, userId, id, category)
			if err != nil {
				return nil, err
			}
		}
	}

	return category, nil
}

func (r *CategoryPlantRepositoryImpl) UpdateCategoryPG(ctx context.Context, category *entities.CategoryPlant) error {
	query := `
			UPDATE categories_plants
			SET category_name = $1,
					category_description = $2,
					updated_at = $3
			WHERE id = $4 AND user_id = $5
	`
	_, err := r.DB.ExecContext(ctx, query,
		category.Name,
		category.Description,
		time.Now(),
		category.Id,
		category.UserId,
	)
	return err
}

func (r *CategoryPlantRepositoryImpl) UpdateCategoryCache(ctx context.Context, category *entities.CategoryPlant) error {
	// Atualiza as chaves Redis para ID e nome
	keyByID := "categories_plants:" + category.UserId + ":" + category.Id
	keyByName := "categories_plants:" + category.UserId + ":" + category.Name

	categoryJSON, err := json.Marshal(category)
	if err != nil {
		return err
	}

	// Atualiza as duas chaves no Redis
	if err := r.RD.Set(ctx, keyByID, categoryJSON, 10*time.Minute).Err(); err != nil {
		return err
	}
	if err := r.RD.Set(ctx, keyByName, categoryJSON, 10*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) Update(ctx context.Context, category *entities.CategoryPlant) error {
	// Atualiza no PostgreSQL
	if err := r.UpdateCategoryPG(ctx, category); err != nil {
		return err
	}

	// Atualiza o cache no Redis
	if err := r.UpdateCategoryCache(ctx, category); err != nil {
		return err
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) DeleteCategoryPG(ctx context.Context, userId, id string) error {
	query := `
			DELETE FROM categories_plants
			WHERE id = $1 AND user_id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, id, userId)
	return err
}

func (r *CategoryPlantRepositoryImpl) DeleteCategoryCache(ctx context.Context, userId, id, name string) error {
	keyByID := "categories_plants:" + userId + ":" + id
	keyByName := "categories_plants:" + userId + ":" + name + ":" + id

	if err := r.RD.Del(ctx, keyByID).Err(); err != nil {
		return err
	}
	if err := r.RD.Del(ctx, keyByName).Err(); err != nil {
		return err
	}

	keyAllCategories := "categories_plants:" + userId
	if err := r.RD.HDel(ctx, keyAllCategories, id).Err(); err != nil {
		return err
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) Delete(ctx context.Context, userId, id string) error {
	category, err := r.FindByIDPG(ctx, userId, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	if err := r.DeleteCategoryPG(ctx, userId, id); err != nil {
		return err
	}

	if err := r.DeleteCategoryCache(ctx, userId, id, category.Name); err != nil {
		return err
	}

	return nil
}
