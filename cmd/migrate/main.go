package main

import (
	"lab1/internal/app/ds"
	"lab1/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("не удалось подключиться к базе данных")
	}

	// 1️⃣ Миграция родительских таблиц без Installations
	if err := db.AutoMigrate(
		&ds.Users{},
		&ds.Software{},
		&ds.SoftwareRequest{},
	); err != nil {
		panic("не удалось выполнить миграцию родительских таблиц")
	}

	// 2️⃣ Миграция InstallationTime
	if err := db.AutoMigrate(&ds.InstallationTime{}); err != nil {
		panic("не удалось выполнить миграцию InstallationTime")
	}

	log.Println("Миграция прошла успешно!")
}
