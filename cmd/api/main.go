package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/tufee/desk-reservation-go/internal/api"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}
}

func init() {
	db, err := sql.Open(os.Getenv("POSTGRES_DB"), os.Getenv("CONNECTION_STRING"))
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Error creating postgres driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("MIGRATION_DIR"),
		os.Getenv("DRIVER"),
		driver,
	)
	if err != nil {
		log.Fatal("Error initializing migration instance:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error executing migrations:", err)
	}
}

func main() {
	router := api.SetupRoutes()

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port 8080")
	server.ListenAndServe()
}
