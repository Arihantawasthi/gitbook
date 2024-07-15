package storage

import (
	"fmt"
	"log"
	"os"
    "database/sql"
    _ "github.com/lib/pq"
)

var DBConn *sql.DB

type PGStorage struct {
	Host string
	User string
	Pass string
	Name string
	Port string
}

func ConnectPGStorage() {
	pgConfig := &PGStorage{
		Host: os.Getenv("DB_HOST"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Name: os.Getenv("DB_NAME"),
		Port: os.Getenv("DB_PORT"),
	}

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgConfig.User,
		pgConfig.Pass,
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.Name,
	)
	conn, err := sql.Open("postgres", dbConnStr)

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DATABASE: %v\n", err)
	}
    if pingErr := conn.Ping(); pingErr != nil {
        log.Fatalf("ERROR PINGING DATABASE: %v\n", pingErr)
    }
    log.Printf("Connection to database established: %v\n", pgConfig.Name)

    DBConn = conn
}
