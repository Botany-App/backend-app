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
	GetUserRoutes := usecases.NewGetUserUseCase(repository)
	useHandlers := handlers.NewUserHandlers(RegisterUserRoutes, LoginUserRoutes, GetUserRoutes)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register/1", useHandlers.RegisterUserHandler)
		r.Post("/register/2", useHandlers.ConfirmEmailHandler)
		r.Post("/register/resendToken", useHandlers.ResendTokenHandler)
		r.Post("/login", useHandlers.LoginUserHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtService))
			r.Get("/user/{ID}", useHandlers.GetByIdUserHandler)
		})
	})

	return r, nil
}
