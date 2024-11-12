package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"

	"github.com/lucasBiazon/botany-back/internal/entities"
)

type UserRepositoryImpl struct {
	DB *sql.DB
	RD *redis.Client
}

func NewUserRepository(db *sql.DB, rd *redis.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DB: db,
		RD: rd,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (ID, name_user, email, password_hash) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`, user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.User, error) {
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

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
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

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET name_user=$1, email=$2, password_user=$3, updated_at=$4 WHERE ID=$5`

	_, err := r.DB.Exec(query, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, email, password string) error {
	query := `UPDATE users SET password_hash=$1 WHERE email=$2`
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(query, string(passwordHash), email)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE ID=$1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) StoreToken(ctx context.Context, email, token string) error {

	err := r.RD.HSet(ctx, email, "token", token).Err()
	if err != nil {
		return err
	}
	err = r.RD.Expire(ctx, email, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) ResendToken(ctx context.Context, email string, token string) error {
	tokenStored, err := r.RD.HGet(ctx, email, "token").Result()
	if err == redis.Nil || tokenStored == "" {
		err = r.RD.HSet(ctx, email, "token", token).Err()
		if err != nil {
			return err
		}
		err = r.RD.Expire(ctx, email, 10*time.Minute).Err()
		if err != nil {
			return err
		}
		return nil
	}
	if err = r.RD.Del(ctx, email).Err(); err != nil {
		return err
	}
	err = r.RD.HSet(ctx, email, "token", token).Err()
	if err != nil {
		return err
	}
	err = r.RD.Expire(ctx, email, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) ActivateAccount(ctx context.Context, email, token string) error {

	tokenStored, err := r.RD.HGet(ctx, email, "token").Result()
	if err == redis.Nil || tokenStored != token {
		return err
	}
	if err = r.RD.Del(ctx, email).Err(); err != nil {
		return err
	}

	query := `UPDATE users SET isActive=true WHERE email=$1`
	_, err = r.DB.Exec(query, email)
	return err
}

func (r *UserRepositoryImpl) Login(ctx context.Context, email, password string) (string, error) {
	query := `SELECT password_hash FROM users WHERE email=$1`
	row := r.DB.QueryRow(query, email)
	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}
	return "token", nil
}
