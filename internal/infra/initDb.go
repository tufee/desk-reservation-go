package infra

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type Db struct {
	Conn *sqlx.DB
}

func InitializeDB() (*Db, error) {
	connection := os.Getenv("CONNECTION_STRING")

	if connection == "" {
		return nil, fmt.Errorf("missing DATABASE_URL in environment")
	}

	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &Db{Conn: db}, nil
}
