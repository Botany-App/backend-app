package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lucasBiazon/botany-back/internal/database"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	database.ConnectDB()
}
