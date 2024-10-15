package repositories

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UserRepositoryPg struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepositoryPg {
	return &UserRepositoryPg{
		DB: db,
	}
}

func (r *UserRepositoryPg) Create(user *entities.User) error {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	_, err := r.DB.Exec(context.Background(), query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPg) GetByID(id uuid.UUID) (*entities.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=$1`

	row := r.DB.QueryRow(context.Background(), query, id)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryPg) GetByEmail(email string) (*entities.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=$1`

	row := r.DB.QueryRow(context.Background(), query, email)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryPg) Update(user *entities.User) error {
	query := `UPDATE users SET name=$1, email=$2, password=$3, updated_at=$4 WHERE id=$5`

	_, err := r.DB.Exec(context.Background(), query, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPg) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id=$1`

	_, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
