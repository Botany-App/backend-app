package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() (*sql.DB, *redis.Client, error) {
	db, err := ConnectPG()
	if err != nil {
		return nil, nil, err
	}

	rd, err := InitRedisClient(context.Background())
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to the database and Redis")
	return db, rd, nil
}

func ConnectPG() (*sql.DB, error) {

	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		log.Fatal("POSTGRES_URL is not set")
	}

	db, err := sql.Open("pgx", postgresURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Verificando a conexão com o banco de dados
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v\n", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create postgres driver: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while running migrations: %v\n", err)
	}

	log.Println("Migrations applied successfully")
	return db, nil
}

func InitRedisClient(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client, nil
}
