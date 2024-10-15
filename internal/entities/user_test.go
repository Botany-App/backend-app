package entities

import (
	"fmt"
	"testing"
)

func TestUserCreation(t *testing.T) {

	t.Run("Create", func(t *testing.T) {

		user, err := NewUser("John Doe", "johnDoe@gmail.com", "password")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if user.Name != "John Doe" {
			t.Errorf("Expected name to be John Doe, got %s", user.Name)
		}

		if user.Email != "johnDoe@gmail.com" {
			t.Errorf("Expected email to be johnDoe@gmail.com, got %s", user.Email)
		}

		if user.Password != "password" {
			t.Errorf("Expected password to be password, got %s", user.Password)
		}
		fmt.Println(user)
	})

	t.Run("Email validation", func(t *testing.T) {

		_, err := NewUser("John Doe", "johnDoe", "password")
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})

	t.Run("Password validation", func(t *testing.T) {

		_, err := NewUser("John Doe", "johnDoe@gmail.com", "")
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})

	t.Run("Name validation", func(t *testing.T) {
		_, err := NewUser("", "johnDoe@gmail.com", "password")
		if err == nil {
			t.Error("Expected an error, got nil")
		}

	})
}
