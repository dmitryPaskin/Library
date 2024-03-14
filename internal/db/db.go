package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"time"
)

type UserDataBase struct {
	DB *sql.DB
}

func New() (UserDataBase, error) {
	if err := godotenv.Load(); err != nil {
		return UserDataBase{nil}, err
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	dataBase, err := sql.Open("postgres", dbURL)

	if err != nil {
		return UserDataBase{nil}, err
	}

	time.Sleep(time.Second * 5)
	if err = dataBase.Ping(); err != nil {
		return UserDataBase{nil}, err
	}

	return UserDataBase{DB: dataBase}, nil
}
