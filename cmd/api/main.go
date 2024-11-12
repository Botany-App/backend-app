package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/lucasBiazon/botany-back/internal/database"
	"github.com/lucasBiazon/botany-back/internal/infra/repositories"
	usecases "github.com/lucasBiazon/botany-back/internal/usecases/user"
	handlers "github.com/lucasBiazon/botany-back/internal/web"
)

func main() {
	// go routines
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	runtime.GOMAXPROCS(1)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Init database and redis
	db, clientRedis, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// // Init user use cases

	repository := repositories.NewUserRepository(db, clientRedis)
	RegisterUserRoutes := usecases.NewRegisterUserUseCase(repository)
	useHandlers := handlers.NewUserHandlers(RegisterUserRoutes)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register/1", useHandlers.RegisterUserHandler)
		r.Post("/register/2", useHandlers.ConfirmEmailHandler)
		r.Post("/register/resendToken", useHandlers.ResendTokenHandler)
	})

	// Start server
	local := os.Getenv("API_PORT")
	if local == "" {
		local = "8081"
	}
	go func() {
		defer wg.Done()
		http.ListenAndServe(fmt.Sprintf(":%s", local), r)
	}()
	fmt.Printf("Server running on port %s\n", local)

	wg.Wait()
}
