package models

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error

	// connect to database
	dsn := os.Getenv("POSTGRES_URL_GORM")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
		os.Exit(1)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	DB.Logger = logger.Default.LogMode(logger.Info)

	// migrate database
	log.Println("Starting Migrations")
	err = DB.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Error migrating database")
		os.Exit(1)
	}

	log.Println("Database Migrated Successfully")
}

func Cleanup() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error getting sqlDB")
		os.Exit(1)
	}

	sqlDB.Close()
}
