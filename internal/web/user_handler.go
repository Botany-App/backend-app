package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
)

type UserHandlers struct {
	RegisterUserUseCase *usecases.RegisterUserUseCase
	LoginUserUseCase    *usecases.LoginUserUseCase
	GetUserUseCase      *usecases.GetUserUseCase
}

func NewUserHandlers(registerUserUseCase *usecases.RegisterUserUseCase, loginUserUseCase *usecases.LoginUserUseCase, getUserUseCase *usecases.GetUserUseCase) *UserHandlers {
	return &UserHandlers{
		RegisterUserUseCase: registerUserUseCase,
		LoginUserUseCase:    loginUserUseCase,
		GetUserUseCase:      getUserUseCase,
	}
}

func (h *UserHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.RegisterUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.RegisterUserUseCase.StartRegistration(context.Background(), input); err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("Erro ao registrar usuário: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("Código enviado para: %s", input.Email))
}

func (h *UserHandlers) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.ConfirmEmailInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.RegisterUserUseCase.ConfirmEmail(context.Background(), input.Email, input.Token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Sprintf("Erro ao confirmar email: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("Email confirmado: %s \n Conta criada com sucesso!", input.Email))

}

func (h *UserHandlers) ResendTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.ResendTokenInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.RegisterUserUseCase.ResendToken(context.Background(), input.Email); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Sprintf("Erro ao reenviar token: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("Código reenviado para: %s", input.Email))
}

func (h *UserHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.LoginUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := h.LoginUserUseCase.Execute(context.Background(), input.Email, input.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Sprintf("Erro ao logar: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("user logged in: %s", token))
}

func (h *UserHandlers) GetByIdUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.GetUserUseCase.GetUserById(context.Background(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(fmt.Sprintf("Erro ao buscar usuário: %s", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
