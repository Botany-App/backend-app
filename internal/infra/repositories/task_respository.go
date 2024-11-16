package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type TaskRepositoryImpl struct {
	DB *sql.DB
	RD *redis.Client
}

func NewTaskRepository(db *sql.DB, rd *redis.Client) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{
		DB: db,
		RD: rd,
	}
}

func (r *TaskRepositoryImpl) Create(task *entities.Task) error {
	query := `INSERT INTO tasks (ID, name_task, description_task, task_date, user_id, task_status) 
          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(query, task.ID, task.Name, task.Description, task.TaskDate, task.UserID, task.TaskStatus)

	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) FindAll(userID string) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:all", userID)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1`
		rows, err := r.DB.Query(query, userID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var task entities.Task
			err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindByID(userID, id string) (*entities.Task, error) {
	key := fmt.Sprintf("task:user:%s:id:%s", userID, id)
	fetch := func() (*entities.Task, error) {
		var task entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1 AND ID = $2`
		err := r.DB.QueryRow(query, userID, id).Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
		if err != nil {
			return nil, err
		}
		return &task, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindAllByName(userID, name string) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:name:%s", userID, name)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1 AND name_task ILIKE '%' || $2 || '%'`
		rows, err := r.DB.Query(query, userID, name)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var task entities.Task
			err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindAllByDate(userID string, date string) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:date:%s", userID, date)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1 AND task_date::date = $2::date`
		rows, err := r.DB.Query(query, userID, date)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var task entities.Task
			err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindAllByStatus(userID, status string) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:status:%s", userID, status)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1 AND task_status = $2`
		rows, err := r.DB.Query(query, userID, status)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var task entities.Task
			err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindTasksNearDeadline(userID string, days int) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:near_deadline:%d", userID, days)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1 AND task_date BETWEEN NOW() AND NOW() + INTERVAL '1 day' * $2`
		rows, err := r.DB.Query(query, userID, days)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var task entities.Task
			err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindTasksFarFromDeadline(userID string, days int) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:far_deadline:%d", userID, days)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks WHERE user_id = $1 AND task_date > NOW() + INTERVAL '1 day' * $2`
		rows, err := r.DB.Query(query, userID, days)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var task entities.Task
			err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TaskDate, &task.UserID, &task.TaskStatus)
			if err != nil {
				return nil, err
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	return getFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) Update(task *entities.Task) error {
	query := `UPDATE tasks SET name_task = $1, description_task = $2, task_date = $3, task_status = $4 WHERE ID = $5`
	_, err := r.DB.Exec(query, task.Name, task.Description, task.TaskDate, task.TaskStatus, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) Delete(id string) error {
	query := `DELETE FROM tasks WHERE ID = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

const cacheDuration = 6 * time.Hour

// Função auxiliar para buscar no cache ou no PostgreSQL
func getFromCache[T any](rd *redis.Client, key string, fetch func() (T, error)) (T, error) {
	var result T
	ctx := context.Background()

	// Tenta buscar o cache
	cachedData, err := rd.Get(ctx, key).Result()
	if err == redis.Nil { // Cache não encontrado
		// Busca no PostgreSQL
		fetchedData, fetchErr := fetch()
		if fetchErr != nil {
			return result, fetchErr
		}

		// Serializa e salva no Redis
		serializedData, err := json.Marshal(fetchedData)
		if err != nil {
			return result, fmt.Errorf("failed to serialize data for caching: %w", err)
		}

		err = rd.Set(ctx, key, serializedData, cacheDuration).Err()
		if err != nil {
			fmt.Printf("failed to cache data: %v\n", err)
		}

		return fetchedData, nil
	} else if err != nil {
		return result, err
	}

	// Desserializa o resultado do cache
	json.Unmarshal([]byte(cachedData), &result)
	return result, nil
}
