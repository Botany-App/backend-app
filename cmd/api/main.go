package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lucasBiazon/botany-back/internal/database"
	"github.com/lucasBiazon/botany-back/internal/infra/repositories"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
	handlers "github.com/lucasBiazon/botany-back/internal/web"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		panic(err)
	}

	repository := repositories.NewUserRepository(db)
	CreateUserUseCase := usecases.NewCreateUserUseCase(repository)
	GetByIdUserUseCase := usecases.NewGetByIdUserUseCase(repository)
	GetByEmailUserUseCase := usecases.NewGetByEmailUserUseCase(repository)

	userHandlers := handlers.NewUserHandler(CreateUserUseCase, GetByIdUserUseCase, GetByEmailUserUseCase)
	r := chi.NewRouter()
	r.Post("/user", userHandlers.CreateUserHandler)
	r.Get("/user", userHandlers.GetUserByIdHandler)
	r.Get("/user/email", userHandlers.GetUserByEmailHandler)

	http.ListenAndServe(":8080", r)
	log.Println("Server running on port 8080")
}
