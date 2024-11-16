package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
)

type UserHandlers struct {
	RegisterUserUseCase             *usecases.RegisterUserUseCase
	LoginUserUseCase                *usecases.LoginUserUseCase
	GetUserByIdUseCase              *usecases.GetUserByIdUseCase
	DeleteUserUseCase               *usecases.DeleteUserUseCase
	UpdateUserUseCase               *usecases.UpdateUserUseCase
	RequestPasswordResetUserUseCase *usecases.RequestPasswordResetUserUseCase
	ResetPasswordUserUseCase        *usecases.ResetPasswordUserUseCase
}

func NewUserHandlers(registerUserUseCase *usecases.RegisterUserUseCase, loginUserUseCase *usecases.LoginUserUseCase,
	getUserByIdUseCase *usecases.GetUserByIdUseCase, deleteUserUseCase *usecases.DeleteUserUseCase, updateUserUseCase *usecases.UpdateUserUseCase,
	requestPasswordResetUserUseCase *usecases.RequestPasswordResetUserUseCase, resetPasswordUserUseCase *usecases.ResetPasswordUserUseCase) *UserHandlers {
	return &UserHandlers{
		RegisterUserUseCase:             registerUserUseCase,
		LoginUserUseCase:                loginUserUseCase,
		GetUserByIdUseCase:              getUserByIdUseCase,
		DeleteUserUseCase:               deleteUserUseCase,
		UpdateUserUseCase:               updateUserUseCase,
		RequestPasswordResetUserUseCase: requestPasswordResetUserUseCase,
		ResetPasswordUserUseCase:        resetPasswordUserUseCase,
	}
}

func (h *UserHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.RegisterUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	if err := h.RegisterUserUseCase.StartRegistration(context.Background(), input); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao registrar usuário", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Código enviado para o email", map[string]string{"email": input.Email})
}

func (h *UserHandlers) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.ConfirmEmailInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	if err := h.RegisterUserUseCase.ConfirmEmail(context.Background(), input); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao confirmar email", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Email confirmado e conta criada com sucesso", nil)
}

func (h *UserHandlers) ResendTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.ResendTokenInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	if err := h.RegisterUserUseCase.ResendToken(context.Background(), input); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao reenviar código", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Código reenviado para o email", map[string]string{"email": input.Email})
}

func (h *UserHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.LoginUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	token, err := h.LoginUserUseCase.Execute(context.Background(), input)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao logar", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Login realizado com sucesso", map[string]string{"token": token})
}

func (h *UserHandlers) GetByIdUserHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	user, err := h.GetUserByIdUseCase.Execute(context.Background(), usecases.GetUserByIdInputDTO{ID: userID})
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao buscar usuário", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Usuário encontrado", user)
}

func (h *UserHandlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	if err := h.DeleteUserUseCase.Execute(context.Background(), usecases.DeleteUserInputDTO{ID: userID}); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao deletar usuário", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Usuário deletado com sucesso", nil)
}

func (h *UserHandlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	userID, err := services.ExtractUserIDFromToken(auth, services.NewJWTService(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		jsonResponse(w, http.StatusUnauthorized, "error", "Token inválido ou expirado", nil)
		return
	}
	var input usecases.UpdateUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	input.ID = userID
	if err := h.UpdateUserUseCase.Execute(context.Background(), input); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao atualizar usuário", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Usuário atualizado com sucesso", nil)
}

func (h *UserHandlers) RequestPasswordResetUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.RequestPasswordResetUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	if err := h.RequestPasswordResetUserUseCase.Execute(context.Background(), input); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao solicitar reset de senha", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Token de redefinição enviado para o e-mail", nil)
}

func (h *UserHandlers) ResetPasswordUserHandler(w http.ResponseWriter, r *http.Request) {
	var input usecases.ResetPasswordUserInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonResponse(w, http.StatusBadRequest, "error", "Erro ao decodificar a requisição", nil)
		return
	}
	if err := h.ResetPasswordUserUseCase.Execute(context.Background(), input); err != nil {
		jsonResponse(w, http.StatusInternalServerError, "error", "Erro ao resetar senha", err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, "success", "Senha resetada com sucesso", nil)
}
