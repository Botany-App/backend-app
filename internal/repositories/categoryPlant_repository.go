package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *CategoryPlantRepositoryImpl) Create(ctx context.Context, category *entities.CategoryPlant) error {
	query := `INSERT INTO categories_plants (ID, name_category, description_category, user_id) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, category.ID, category.Name, category.Description, category.UserID)
	if err != nil {
		return errors.New("error creating category plant")
	}
	return nil
}

func (r *CategoryPlantRepositoryImpl) FindAll(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.CategoryPlant, error) {
	cacheKey := fmt.Sprintf("categories_plant_user_%s", userID)

	fetchFromDB := func() ([]*entities.CategoryPlant, error) {
		query := `SELECT ID, name_category, description_category, created_at, updated_at FROM categories_plants WHERE user_id = $1 LIMIT $2 OFFSET $3`
		rows, err := r.DB.Query(query, userID, limit, offset)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var categories []*entities.CategoryPlant
		for rows.Next() {
			var category entities.CategoryPlant
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
	categoriesResult, err := GetFromCache(r.RD, cacheKey, fetchFromDB)
	if err != nil {
		return nil, err
	}
	if len(categoriesResult) == 0 {
		return nil, errors.New("categories not found")
	}

	return categoriesResult, nil
}

func (r *CategoryPlantRepositoryImpl) FindByID(ctx context.Context, userID, id uuid.UUID, limit, offset int) (*entities.CategoryPlant, error) {
	cacheKey := fmt.Sprintf("categories_plant_user_%s_id_%s", userID, id)

	fetchFromDB := func() (*entities.CategoryPlant, error) {
		query := `SELECT ID, name_category, description_category, created_at, updated_at FROM categories_plants WHERE user_id = $1 AND ID = $2 LIMIT $3 OFFSET $4`
		row := r.DB.QueryRow(query, userID, id, limit, offset)

		var category entities.CategoryPlant
		err := row.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		category.UserID = userID
		return &category, nil
	}
	categoryResult, err := GetFromCache(r.RD, cacheKey, fetchFromDB)
	if err != nil {
		return nil, err
	}
	if categoryResult == nil {
		return nil, errors.New("category not found")
	}

	return categoryResult, nil
}

func (r *CategoryPlantRepositoryImpl) FindByName(ctx context.Context, userID uuid.UUID, name string, limit, offset int) ([]*entities.CategoryPlant, error) {
	cacheKey := fmt.Sprintf("categories_plant_user_%s_name_%s", userID, name)

	fetchFromDB := func() ([]*entities.CategoryPlant, error) {
		query := `SELECT ID, name_category, description_category, created_at, updated_at FROM categories_plants WHERE user_id = $1 AND name_category ILIKE '%' || $2 || '%'`
		rows, err := r.DB.Query(query, userID, name)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		var categories []*entities.CategoryPlant
		for rows.Next() {
			var category entities.CategoryPlant
			err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
			if err != nil {
				return nil, err
			}
			category.UserID = userID
			categories = append(categories, &category)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}

		return categories, nil
	}
	categoriesResult, err := GetFromCache(r.RD, cacheKey, fetchFromDB)
	if err != nil {
		return nil, err
	}
	if categoriesResult == nil {
		return nil, errors.New("category not found")
	}

	if len(categoriesResult) == 0 {
		return nil, errors.New("categories not found")
	}

	return categoriesResult, nil
}

func (r *CategoryPlantRepositoryImpl) Update(ctx context.Context, category *entities.CategoryPlant) error {
	query := `UPDATE categories_plants SET name_category = $1, description_category = $2 WHERE ID = $3 AND user_id = $4`
	_, err := r.DB.Exec(query, category.Name, category.Description, category.ID, category.UserID)
	if err != nil {
		return errors.New("error updating category plant")
	}

	cacheKey := fmt.Sprintf("categories_plant_user_%s_id_%s", category.UserID, category.ID)
	err = r.RD.Del(ctx, cacheKey).Err()
	if err != nil {
		return errors.New("error deleting cache")
	}
	return nil
}

func (r *CategoryPlantRepositoryImpl) Delete(ctx context.Context, userID, id uuid.UUID) error {
	query := `DELETE FROM categories_plants WHERE ID = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return errors.New("error deleting category plant")
	}

	cacheKey := fmt.Sprintf("categories_plant_user_%s_id_%s", userID, id)
	err = r.RD.Del(ctx, cacheKey).Err()
	if err != nil {
		return errors.New("error deleting cache")
	}
	return nil
}
