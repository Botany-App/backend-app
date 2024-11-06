package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/lucasBiazon/botany-back/internal/entities"
	"github.com/lucasBiazon/botany-back/internal/utils"
)

type UserRepositoryPg struct {
	DB *sql.DB
	RD *redis.Client
}

func NewUserRepository(db *sql.DB, rd *redis.Client) *UserRepositoryPg {
	return &UserRepositoryPg{
		DB: db,
		RD: rd,
	}
}

func (r *UserRepositoryPg) Create(user *entities.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (ID, name_user, email, password_hash) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`, user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryPg) GetByID(id string) (*entities.User, error) {
	query := `SELECT ID, name_user, email, password_hash, created_at, updated_at FROM users WHERE id=$1`

	row := r.DB.QueryRow(query, id)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryPg) GetByEmail(email string) (*entities.User, error) {
	query := `SELECT ID, name_user, email, password_hash, created_at, updated_at FROM users WHERE email=$1`

	row := r.DB.QueryRow(query, email)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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

func (r *UserRepositoryPg) HostCode(user *entities.User) error {
	var ctx = context.Background()
	code, err := utils.GenerateCode()
	if err != nil {
		return err
	}

	utils.SendEmail(user.Email, code)
	cacheKey := fmt.Sprintf("user:%s", user.Email)
	err = r.RD.HSet(ctx, cacheKey, "id", user.ID.String(), "name", user.Name, "email", user.Email, "password", user.Password, "code", code).Err()
	if err != nil {
		return err
	}

	err = r.RD.Expire(ctx, cacheKey, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryPg) VerifyEmail(email, codeInput string) (*entities.User, error) {
	var ctx = context.Background()

	code, err := r.RD.HGet(ctx, email, "code").Result()
	if err == redis.Nil {
		return nil, errors.New("code not found")
	}
	if err != nil {
		return nil, err
	}
	if code != codeInput {
		return nil, errors.New("invalid code")
	}

	name, err := r.RD.HGet(ctx, email, "name").Result()
	if err != nil {
		return nil, err
	}
	password, err := r.RD.HGet(ctx, email, "password").Result()
	if err != nil {
		return nil, err
	}

	emailUser, err := r.RD.HGet(ctx, email, "email").Result()
	if err != nil {
		return nil, err
	}
	user, err := entities.NewUser(name, emailUser, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPg) Login(email, password string) (*entities.User, error) {
	query := `SELECT ID, name_user, email, password_hash, created_at, updated_at FROM users WHERE email=$1 AND password_hash=$2`

	row := r.DB.QueryRow(query, email, password)
	user := &entities.User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
