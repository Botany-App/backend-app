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

	key := "categories_plants:" + categoryPlant.UserId
	categoryData, err := json.Marshal(categoryPlant)
	if err != nil {
		return "", err
	}

	_, err = r.RD.HSet(ctx, key, categoryPlant.Id, categoryData).Result()
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

	return r.RD.Expire(ctx, key, 20*time.Minute).Err() // Configura TTL
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

func (r *CategoryPlantRepositoryImpl) SetCategoryByNameRD(ctx context.Context, userId string, categories []*entities.CategoryPlant) error {
	for _, category := range categories {
		key := "categories_plants:" + userId + ":" + category.Name + ":" + category.Id

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

func (r *CategoryPlantRepositoryImpl) FindByIDPG(ctx context.Context, userId, id string) (*entities.CategoryPlant, error) {
	query := `
		SELECT id, category_name, category_description, user_id, created_at, updated_at
		FROM categories_plants
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
		return nil, fmt.Errorf("erro ao escanear categoria: %w", err)
	}

	return &category, nil
}

func (r *CategoryPlantRepositoryImpl) FindByIDRD(ctx context.Context, userId, id string) (*entities.CategoryPlant, error) {
	// Define a chave para busca no Redis
	key := "categories_plants:" + userId + ":" + id

	// Tenta recuperar os dados do Redis
	categoryJSON, err := r.RD.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria no Redis: %w", err)
	}

	// Desserializa o JSON para a entidade
	var category entities.CategoryPlant
	if err := json.Unmarshal([]byte(categoryJSON), &category); err != nil {
		return nil, fmt.Errorf("erro ao desserializar categoria: %w", err)
	}

	return &category, nil
}

func (r *CategoryPlantRepositoryImpl) SetCategoryByIDRD(ctx context.Context, userId, id string, category *entities.CategoryPlant) error {
	// Define a chave e serializa a entidade para JSON
	key := "categories_plants:" + userId + ":" + id

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

func (r *CategoryPlantRepositoryImpl) FindById(ctx context.Context, userId, id string) (*entities.CategoryPlant, error) {
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

func (r *CategoryPlantRepositoryImpl) UpdateCategoryPG(ctx context.Context, category *entities.CategoryPlant) error {
	query := `
		UPDATE categories_plants
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

func (r *CategoryPlantRepositoryImpl) UpdateCategoryCache(ctx context.Context, category *entities.CategoryPlant) error {
	keyById := "categories_plants:" + category.UserId + ":" + category.Id
	keyByName := "categories_plants:" + category.UserId + ":" + category.Name + ":" + category.Id
	keyAllCategories := "categories_plants:" + category.UserId

	// Remove os caches relacionados
	if err := r.RD.Del(ctx, keyById, keyByName).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache por ID e Nome: %w", err)
	}
	if err := r.RD.HDel(ctx, keyAllCategories, category.Id).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache de categorias gerais: %w", err)
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) Update(ctx context.Context, category *entities.CategoryPlant) error {
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

func (r *CategoryPlantRepositoryImpl) DeleteCategoryPG(ctx context.Context, userId, id string) error {
	query := `
		DELETE FROM categories_plants
		WHERE id = $1 AND user_id = $2
	`
	_, err := r.DB.ExecContext(ctx, query, id, userId)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria no PostgreSQL: %w", err)
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) DeleteCategoryCache(ctx context.Context, userId, id, name string) error {
	keyByID := "categories_plants:" + userId + ":" + id
	keyByName := "categories_plants:" + userId + ":" + name + ":" + id
	keyAllCategories := "categories_plants:" + userId

	// Remove os caches relacionados
	if err := r.RD.Del(ctx, keyByID, keyByName).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache por ID e Nome: %w", err)
	}
	if err := r.RD.HDel(ctx, keyAllCategories, id).Err(); err != nil {
		return fmt.Errorf("erro ao remover cache de categorias gerais: %w", err)
	}

	return nil
}

func (r *CategoryPlantRepositoryImpl) Delete(ctx context.Context, userId, id string) error {
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
