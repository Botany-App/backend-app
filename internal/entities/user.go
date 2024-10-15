package entities

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
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

	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
