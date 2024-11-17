package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/pkg/errors"
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

func (r *CategoryTaskRepositoryImpl) GetAll(userID string) ([]entities.CategoryTask, error) {
	log.Print(userID)
	key := fmt.Sprintf("category_tasks:user:%s:all", userID)
	fetch := func() ([]entities.CategoryTask, error) {
		var categoryTasks []entities.CategoryTask
		query := `SELECT category_id FROM categories_tasks_users WHERE user_id = $1`
		var categoryIDs []string
		rows, err := r.DB.Query(query, userID)
		log.Print(rows)
		if err != nil {
			return nil, errors.New("error fetching data1")
		}
		for rows.Next() {
			var categoryID string
			err := rows.Scan(&categoryID)
			if err != nil {
				return nil, errors.New("error fetching data2")
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		log.Print(categoryIDs)
		query = `SELECT * FROM categories_tasks WHERE id = $1`
		for _, categoryID := range categoryIDs {
			var categoryTask entities.CategoryTask
			err := r.DB.QueryRow(query, categoryID).Scan(&categoryTask.ID, &categoryTask.Name, &categoryTask.Description, &categoryTask.CreatedAt, &categoryTask.UpdatedAt)
			if err != nil {
				return nil, errors.New("error fetching data3")
			}
			categoryTasks = append(categoryTasks, categoryTask)
		}
		log.Println(categoryTasks)
		return categoryTasks, nil
	}
	return GetFromCache(r.RD, key, fetch)

}

func (r *CategoryTaskRepositoryImpl) GetByName(userID string, name string) ([]entities.CategoryTask, error) {
	key := fmt.Sprintf("category_tasks:user:%s:name:%s", userID, name)
	fetch := func() ([]entities.CategoryTask, error) {
		var categoryTasks []entities.CategoryTask
		query := `SELECT category_id FROM categories_tasks_users WHERE user_id = $1`
		var categoryIDs []string
		rows, err := r.DB.Query(query, userID)
		if err != nil {
			return nil, errors.New("error fetching data")
		}
		for rows.Next() {
			var categoryID string
			err := rows.Scan(&categoryID)
			if err != nil {
				return nil, errors.New("error fetching data")
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		query = `SELECT * FROM categories_tasks WHERE id = $1 AND name ILIKE '%' || $2 || '%'`
		for _, categoryID := range categoryIDs {
			var categoryTask entities.CategoryTask
			err := r.DB.QueryRow(query, categoryID, name).Scan(&categoryTask.ID, &categoryTask.Name, &categoryTask.Description, &categoryTask.CreatedAt, &categoryTask.UpdatedAt)
			if err != nil {
				return nil, errors.New("error fetching data")
			}
			categoryTasks = append(categoryTasks, categoryTask)
		}
		return categoryTasks, nil
	}

	return GetFromCache(r.RD, key, fetch)
}

func (r *CategoryTaskRepositoryImpl) GetByID(userID, id string) ([]entities.CategoryTask, error) {
	key := fmt.Sprintf("category_tasks:user:%s:id:%s", userID, id)
	fetch := func() ([]entities.CategoryTask, error) {
		var categoryTasks []entities.CategoryTask
		query := `SELECT category_id FROM categories_tasks_users WHERE user_id = $1`
		var categoryIDs []string
		rows, err := r.DB.Query(query, userID)
		if err != nil {
			return nil, errors.New("error fetching data")
		}
		for rows.Next() {
			var categoryID string
			err := rows.Scan(&categoryID)
			if err != nil {
				return nil, errors.New("error fetching data")
			}
			categoryIDs = append(categoryIDs, categoryID)
		}
		query = `SELECT * FROM categories_tasks WHERE id = $1`
		for _, categoryID := range categoryIDs {
			var categoryTask entities.CategoryTask
			err := r.DB.QueryRow(query, categoryID).Scan(&categoryTask.ID, &categoryTask.Name, &categoryTask.Description, &categoryTask.CreatedAt, &categoryTask.UpdatedAt)
			if err != nil {
				return nil, errors.New("error fetching data")
			}
			categoryTasks = append(categoryTasks, categoryTask)
		}
		return categoryTasks, nil
	}

	return GetFromCache(r.RD, key, fetch)
}

func (r *CategoryTaskRepositoryImpl) Create(userID string, category *entities.CategoryTask) error {
	query := `INSERT INTO categories_tasks (ID, name_category , description_category) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {

		return errors.New("primeira error creating category task")
	}

	query = `INSERT INTO categories_tasks_users (ID, category_id, user_id) VALUES ($1, $2, $3)`
	_, err = r.DB.Exec(query, uuid.New(), category.ID, userID)
	if err != nil {
		return errors.New("segunda error creating category task")
	}

	return nil
}

func (r *CategoryTaskRepositoryImpl) Update(userID string, category *entities.CategoryTask) error {
	query := `UPDATE categories_tasks SET name = $1, description = $2 WHERE id = $3`
	_, err := r.DB.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return errors.New("error updating category task")
	}
	return nil
}

func (r *CategoryTaskRepositoryImpl) Delete(userID, id string) error {
	query := `DELETE FROM categories_tasks WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return errors.New("error deleting category task")
	}
	return nil
}
