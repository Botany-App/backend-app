package entities

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	HostCode(*User) error
	VerifyEmail(email string, codeInput string) (*User, error)
	Login(email string, password string) (*User, error)
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
