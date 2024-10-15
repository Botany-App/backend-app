package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"github.com/lucasBiazon/botany-back/internal/entities"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	conect := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	conn, err := pgx.Connect(conect)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	log.Println("Server is running on port 8080")
	go http.ListenAndServe(":8080", r)

	user, err := entities.NewUser("John Doe", "jhon@gmail.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
