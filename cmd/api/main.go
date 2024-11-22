package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
	"github.com/lucasBiazon/botany-back/internal/database"
	"github.com/lucasBiazon/botany-back/internal/routes"
	services "github.com/lucasBiazon/botany-back/internal/service"
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
		log.Panic(err)
	}

	// // Init user use cases
	secretKey := os.Getenv("JWT_SECRET_KEY")
	jwtService := services.NewJWTService(secretKey)
	r, err := routes.InitializeRoutes(db, clientRedis, jwtService)
	if err != nil {
		panic(err)
	}

	// Start server
	local := os.Getenv("API_PORT")
	if local == "" {
		local = "8080"
	}
	go func() {
		defer wg.Done()
		for {
			http.ListenAndServe(fmt.Sprintf(":%s", local), r)
		}
	}()
	log.Printf("Server running on port %s\n", local)

	wg.Wait()
}
