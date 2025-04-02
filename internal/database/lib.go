package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func Init() {
	db, err = gorm.Open(postgres.Open(os.Getenv("POSTGRES_DSN")), &gorm.Config{})

	if err != nil {
		log.Fatal("Database failed to initialize: %w", err)
	}

	db.Exec("CREATE SCHEMA IF NOT EXISTS firstweb")
	db.Exec("SET search_path TO firstweb")

	db.AutoMigrate(&User{})
}
