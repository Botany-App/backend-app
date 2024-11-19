package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

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

func (r *CategoryTaskRepositoryImpl) FindAll(ctx context.Context, userID string) ([]*entities.CategoryTask, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for user_id: %w", err)
	}
	cacheKey := fmt.Sprintf("categories_task_user_%s", userID)

	fetchFromDB := func() ([]*entities.CategoryTask, error) {
		query := `SELECT ID, name_category, description_category, created_at, updated_at FROM categories_tasks WHERE user_id = $1`
		rows, err := r.DB.Query(query, userUUID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var categories []*entities.CategoryTask
		for rows.Next() {
			var category entities.CategoryTask
			err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
			if err != nil {
				return nil, err
			}
			category.UserID = userID
			categories = append(categories, &category)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return categories, nil
	}

	categories, err := GetFromCache(r.RD, cacheKey, fetchFromDB)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, errors.New("nenhuma categoria de tarefa encontrada")
	}
	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) FindByName(ctx context.Context, userID string, name string) ([]*entities.CategoryTask, error) {
	cacheKey := fmt.Sprintf("categories_task_user_%s_name_%s", userID, name)

	fetchFromDB := func() ([]*entities.CategoryTask, error) {
		query := `SELECT ID, name_category, description_category, created_at, updated_at FROM categories_tasks WHERE user_id = $1 AND name_category ILIKE '%' || $2 || '%'`
		rows, err := r.DB.Query(query, userID, name)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var categories []*entities.CategoryTask
		for rows.Next() {
			var category entities.CategoryTask
			err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
			if err != nil {
				return nil, err
			}
			category.UserID = userID
			categories = append(categories, &category)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		return categories, nil
	}

	categories, err := GetFromCache(r.RD, cacheKey, fetchFromDB)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, errors.New("nenhuma categoria de tarefa encontrada")
	}

	return categories, nil
}

func (r *CategoryTaskRepositoryImpl) FindByID(ctx context.Context, userID, id string) (*entities.CategoryTask, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for user_id: %w", err)
	}

	taskUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID format for task_id: %w", err)
	}

	cacheKey := fmt.Sprintf("categories_task_user_%s_id_%s", userID, id)

	fetchFromDB := func() (*entities.CategoryTask, error) {
		query := `SELECT ID, name_category, description_category, created_at, updated_at 
				  FROM categories_tasks 
				  WHERE user_id = $1 AND ID = $2`
		row := r.DB.QueryRow(query, userUUID, taskUUID)

		var category entities.CategoryTask
		err := row.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("categoria de tarefa não encontrada")
			}
			return nil, err
		}
		category.UserID = userID
		return &category, nil
	}

	category, err := GetFromCache[*entities.CategoryTask](r.RD, cacheKey, fetchFromDB)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, fmt.Errorf("categoria de tarefa não encontrada")
	}

	return category, nil
}

func (r *CategoryTaskRepositoryImpl) Create(ctx context.Context, category *entities.CategoryTask) error {
	query := `
			INSERT INTO categories_tasks (ID, user_id, name_category, description_category)
			VALUES ($1, $2, $3, $4)
	`

	userUUID, err := uuid.Parse(category.UserID)
	if err != nil {
		return fmt.Errorf("invalid UUID format for user_id: %w", err)
	}

	_, err = r.DB.ExecContext(ctx, query, category.ID, userUUID, category.Name, category.Description)
	if err != nil {
		return fmt.Errorf("failed to insert category task: %w", err)
	}

	return nil
}
func (r *CategoryTaskRepositoryImpl) Update(ctx context.Context, category *entities.CategoryTask) error {
	// Parse do UUID do usuário
	userUUID, err := uuid.Parse(category.UserID)
	if err != nil {
		return errors.New("invalid UUID format for user_id")
	}

	// Query de atualização no banco de dados
	query := `UPDATE categories_tasks 
			  SET name_category = $1, description_category = $2, updated_at = $3 
			  WHERE ID = $4 AND user_id = $5`
	_, err = r.DB.Exec(query, category.Name, category.Description, category.UpdatedAt, category.ID, userUUID)
	if err != nil {
		return err
	}

	// Remover o cache após a atualização
	cacheKey := fmt.Sprintf("categories_task_user_%s_id_%s", category.UserID, category.ID)
	err = r.RD.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("Erro ao remover cache para chave %s: %v", cacheKey, err)
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) Delete(ctx context.Context, userID, id string) error {
	// Query de exclusão no PostgreSQL
	query := `DELETE FROM categories_tasks WHERE ID = $1 AND user_id = $2`
	_, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	// Montagem da chave do cache
	cacheKey := fmt.Sprintf("categories_task_user_%s_id_%s", userID, id)
	log.Printf("Removing cache with key: %s", cacheKey)

	// Remoção do cache no Redis
	err = r.RD.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("Failed to remove cache: %v", err)
		return fmt.Errorf("failed to remove cache: %w", err)
	}

	log.Println("Cache successfully removed")
	return nil
}
