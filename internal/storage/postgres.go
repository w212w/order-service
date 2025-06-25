package storage

import (
	"database/sql"
	"fmt"
	"log"
	"order-service/config"

	_ "github.com/lib/pq"
)

func ConnectDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("DB is not reachable: %v", err)
	}
	log.Println("Connected to the database")
	return db
}
