package handlers

import (
	"encoding/json"
	"net/http"

	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
)

type UserHandlers struct {
	CreateUserUseCase     *usecases.CreateUserUseCase
	GetByEmailUserUseCase *usecases.GetByEmailUserUseCase
	GetByIdUserUseCase    *usecases.GetByIdUserUseCase
}

func NewUserHandler(createUserUseCase *usecases.CreateUserUseCase, getByIdUserUseCase *usecases.GetByIdUserUseCase, getByEmailUserUseCase *usecases.GetByEmailUserUseCase) *UserHandlers {
	return &UserHandlers{
		CreateUserUseCase:     createUserUseCase,
		GetByEmailUserUseCase: getByEmailUserUseCase,
		GetByIdUserUseCase:    getByIdUserUseCase,
	}
}

func (u *UserHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.CreateUserInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := u.CreateUserUseCase.Execute(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

func (u *UserHandlers) GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var email usecases.GetByEmailInputDTO
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := u.GetByEmailUserUseCase.Execute(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (u *UserHandlers) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	var id usecases.GetByIdInputDTO
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := u.GetByIdUserUseCase.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
