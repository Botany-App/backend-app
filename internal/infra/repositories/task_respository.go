package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (r *TaskRepositoryImpl) AddCategory(taskID string, categoryID string) error {
	query := `INSERT INTO categories_and_tasks (ID, task_id, category_id) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, uuid.New(), taskID, categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) AddGardenPlanta(taskID string, gardenPlantaID string) error {
	query := `INSERT INTO task_plants_garden (ID, task_id, plant_garden_id) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, uuid.New(), taskID, gardenPlantaID)
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

	return GetFromCache(r.RD, key, fetch)
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

	return GetFromCache(r.RD, key, fetch)
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

	return GetFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindAllByDate(userID string, date time.Time) ([]entities.Task, error) {
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

	return GetFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindAllByStatus(userID string, status entities.TaskStatusEnum) ([]entities.Task, error) {
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

	return GetFromCache(r.RD, key, fetch)
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

	return GetFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) FindAllByCategory(userID, categoryID string) ([]entities.Task, error) {
	key := fmt.Sprintf("tasks:user:%s:category:%s", userID, categoryID)
	fetch := func() ([]entities.Task, error) {
		var tasks []entities.Task
		query := `SELECT * FROM tasks t
		INNER JOIN categories_and_tasks ct ON t.ID = ct.task_id
		WHERE ct.category_id = $1 AND t.user_id = $2`
		rows, err := r.DB.Query(query, categoryID, userID)
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

	return GetFromCache(r.RD, key, fetch)
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

	return GetFromCache(r.RD, key, fetch)
}

func (r *TaskRepositoryImpl) Update(task *entities.Task) error {

	query := `UPDATE tasks SET name_task = $1, description_task = $2, task_date = $3, task_status = $4 WHERE ID = $5`
	_, err := r.DB.Exec(query, task.Name, task.Description, task.TaskDate, task.TaskStatus, task.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) UpdateTaskCategory(taskID, categoryID string) error {
	query := `UPDATE categories_and_tasks SET category_id = $1 WHERE task_id = $2`
	_, err := r.DB.Exec(query, categoryID, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) UpdateTaskGardenPlanta(taskID, gardenPlantaID string) error {
	query := `UPDATE task_plants_garden SET plant_garden_id = $1 WHERE task_id = $2`
	_, err := r.DB.Exec(query, gardenPlantaID, taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) Delete(userID string, id string) error {
	query := `DELETE FROM tasks WHERE ID = $1 and user_id = $2`
	_, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}
	return nil
}
