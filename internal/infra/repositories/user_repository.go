package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (r *UserRepositoryImpl) GetByID(ctx context.Context, ID string) (*entities.User, error) {
	query := `SELECT ID, name_user, email, password_hash, created_at, updated_at FROM users WHERE id=$1`

	row := r.DB.QueryRow(query, ID)
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
	query := `UPDATE users SET name_user=$1, email=$2  WHERE ID=$3`

	_, err := r.DB.Exec(query, user.Name, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) UpdatePassword(ctx context.Context, ID uuid.UUID, password string) error {
	query := `UPDATE users SET password_hash=$1 WHERE ID=$2`
	_, err := r.DB.Exec(query, password, ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, ID string) error {
	query := `DELETE FROM users WHERE ID=$1`
	_, err := r.DB.Exec(query, ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) StoreRevokedTokenPassword(ctx context.Context, token string) error {
	err := r.RD.Set(ctx, token, "revoked", time.Minute*10).Err()
	if err != nil {
		return errors.New("failed to store revoked token")
	}
	return nil
}

func (r *UserRepositoryImpl) IsTokenRevokedPassword(ctx context.Context, token string) bool {
	val, err := r.RD.Get(ctx, token).Result()
	if err == redis.Nil {
		return false
	}
	return err == nil && val == "revoked"
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

func (r *UserRepositoryImpl) ResendToken(ctx context.Context, email string, token string) (string, error) {
	tokenStored, err := r.RD.HGet(ctx, email, "token").Result()

	if err == redis.Nil || tokenStored == "" {
		err = r.RD.HSet(ctx, email, "token", token).Err()
		if err != nil {
			return "", err
		}
		err = r.RD.Expire(ctx, email, 10*time.Minute).Err()
		if err != nil {
			return "", err
		}
		return token, nil
	}
	if err = r.RD.Del(ctx, email).Err(); err != nil {
		return "", err
	}
	err = r.RD.HSet(ctx, email, "token", token).Err()
	if err != nil {
		return "", err
	}
	err = r.RD.Expire(ctx, email, 10*time.Minute).Err()
	if err != nil {
		return "", err
	}
	return tokenStored, nil
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
	query := `SELECT ID, password_hash FROM users WHERE email=$1`
	row := r.DB.QueryRow(query, email)
	var passwordHash, ID string
	err := row.Scan(&ID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "not found", errors.New("user not found")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "invalid password", errors.New("invalid password")
	}
	return ID, nil
}
