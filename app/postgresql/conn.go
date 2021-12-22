package postgresql

import (
	"database/sql"
	"fmt"
	"online-store-go/pkg/logger_utils"
	"os"

	"github.com/joho/godotenv"
)

const (
	postgres_users_username = "postgres_users_username"
	postgres_users_password = "postgres_users_password"
	postgres_users_host     = "postgres_users_host"
	postgres_users_database = "postgres_users_database"
)

var (
	Client *sql.DB
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv(postgres_users_username),
		os.Getenv(postgres_users_password),
		os.Getenv(postgres_users_host),
		os.Getenv(postgres_users_database),
	)

	Client, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	logger_utils.Info("Database succesfully configured")
}
