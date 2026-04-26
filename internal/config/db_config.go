package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func StartDatabase() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	pgsqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", pgsqlInfo)
	if err != nil {
		log.Fatal("cannot connect to database: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	log.Info("database connected")

	return db
}
