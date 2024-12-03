package entities

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	IsActive  bool      `json:"is_active"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	StoreToken(ctx context.Context, email, token string) error
	ResendToken(ctx context.Context, email string, token string) (string, error)
	ActivateAccount(ctx context.Context, email, token string) error
	Login(ctx context.Context, email, password string) (string, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	StoreRevokedTokenPassword(ctx context.Context, token string) error
	IsTokenRevokedPassword(ctx context.Context, token string) bool
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" {
		return nil, errors.New("expected name")
	}

	if email == "" {
		return nil, errors.New("expected email")
	}

	if password == "" {
		return nil, errors.New("expected password")
	}

	if !isEmailValid(email) {
		return nil, errors.New("invalid email")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:        uuid.New(),
		Name:      name,
		Email:     email,
		IsActive:  false,
		Password:  string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
