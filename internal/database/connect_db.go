package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import necessário para usar arquivos locais
	"github.com/joho/godotenv"
)

func ConnectDB() {
	godotenv.Load(".env")
	postgresURL := os.Getenv("POSTGRES_URL")

	if postgresURL == "" {
		log.Fatal("POSTGRES_URL is not set")
	}

	db, err := sql.Open("pgx", postgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// Configurando o driver postgres para o migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create postgres driver: %v\n", err)
	}

	// Usando o caminho correto para as migrações
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations", // Caminho das migrações
		"postgres", driver)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	// Subir as migrações
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while running migrations: %v\n", err)
	}

	log.Println("Migrations applied successfully")

}
