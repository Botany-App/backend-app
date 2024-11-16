package routes

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/lucasBiazon/botany-back/internal/infra/repositories"
	"github.com/lucasBiazon/botany-back/internal/middleware"
	services "github.com/lucasBiazon/botany-back/internal/service"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
	handlers "github.com/lucasBiazon/botany-back/internal/web"
)

func InitializeRoutes(db *sql.DB, clientRedis *redis.Client, jwtService services.JWTService) (*chi.Mux, error) {

	repository := repositories.NewUserRepository(db, clientRedis)
	RegisterUserRoutes := usecases.NewRegisterUserUseCase(repository)
	LoginUserRoutes := usecases.NewLoginUserUseCase(repository)
	GetUserRoutes := usecases.NewGetUserByIdUseCase(repository)
	DeleteUserRoutes := usecases.NewDeleteUserUseCase(repository)
	UpdateUserRoutes := usecases.NewUpdateUserUseCase(repository)
	RequestPasswordResetUserRoutes := usecases.NewRequestPasswordResetUseCase(repository, jwtService)
	ResetPasswordUserRoutes := usecases.NewResetPasswordUserUseCase(repository, jwtService)

	userHandlers := handlers.NewUserHandlers(
		RegisterUserRoutes,
		LoginUserRoutes,
		GetUserRoutes,
		DeleteUserRoutes,
		UpdateUserRoutes,
		RequestPasswordResetUserRoutes,
		ResetPasswordUserRoutes,
	)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", userHandlers.RegisterUserHandler)
		r.Post("/register/confirm", userHandlers.ConfirmEmailHandler)
		r.Post("/register/resend-token", userHandlers.ResendTokenHandler)
		r.Post("/login", userHandlers.LoginUserHandler)
		r.Post("/password-reset/request", userHandlers.RequestPasswordResetUserHandler)
		r.Post("/password-reset", userHandlers.ResetPasswordUserHandler)
	})

	// Rotas protegidas por autenticação
	r.Route("/api/v1/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(jwtService))
		r.Get("/", userHandlers.GetByIdUserHandler)
		r.Delete("/", userHandlers.DeleteUserHandler)
		r.Put("/", userHandlers.UpdateUserHandler)
	})

	return r, nil
}
