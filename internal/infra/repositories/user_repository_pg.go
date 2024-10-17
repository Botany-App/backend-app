package repositories

import (
	"database/sql"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UserRepositoryPg struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryPg {
	return &UserRepositoryPg{
		DB: db,
	}
}

func (r *UserRepositoryPg) Create(user *entities.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (ID, name_user, email, password_user, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPg) GetByID(id string) (*entities.User, error) {
	query := `SELECT ID, name_user, email, password_user, created_at, updated_at FROM users WHERE id=$1`

	row := r.DB.QueryRow(query, id)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Se não encontrar, retorna nil em vez de erro
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryPg) GetByEmail(email string) (*entities.User, error) {
	query := `SELECT ID, name_user, email, password_user, created_at, updated_at FROM users WHERE email=$1`

	row := r.DB.QueryRow(query, email)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Se não encontrar, retorna nil em vez de erro
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryPg) Update(user *entities.User) error {
	query := `UPDATE users SET name_user=$1, email=$2, password_user=$3, updated_at=$4 WHERE ID=$5`

	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPg) Delete(id string) error {
	query := `DELETE FROM users WHERE ID=$1`

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
