package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/lucasBiazon/botany-back/internal/database"
	"github.com/lucasBiazon/botany-back/internal/infra/repositories"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
	handlers "github.com/lucasBiazon/botany-back/internal/web"
)

func main() {
	db, clientRedis, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	repository := repositories.NewUserRepository(db, clientRedis)
	RegisterUserUseCase := usecases.NewRegisterUserUseCase(repository)
	GetByIdUserUseCase := usecases.NewGetByIdUserUseCase(repository)
	GetByEmailUserUseCase := usecases.NewGetByEmailUserUseCase(repository)
	CreateUserUseCase := usecases.NewCreateUserUseCase(repository)

	userHandlers := handlers.NewUserHandler(CreateUserUseCase, GetByIdUserUseCase, GetByEmailUserUseCase, RegisterUserUseCase)
	r := chi.NewRouter()
	r.Post("/user/register", userHandlers.RegisterUserHandler)
	r.Post("/user/register/create", userHandlers.CreateUserHandler)
	r.Get("/user", userHandlers.GetUserByIdHandler)
	r.Get("/user/email", userHandlers.GetUserByEmailHandler)

	http.ListenAndServe(":8081", r)
	log.Println("Server running on port 8080")

}
