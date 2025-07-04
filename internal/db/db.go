package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"admin",
		"postgres_db_password",
		"192.168.68.82",
		"5432",
		"livestreams",
	)

	var err error
	Pool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Connected to PostgreSQL")
}