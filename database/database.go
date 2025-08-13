package database

import (
	"fmt"
	"log"
	"os"

	"appFashion/backend/models" // Pastikan ini sesuai dengan nama modul Anda

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	var dsn string

	// Render menyediakan satu DATABASE_URL, sementara lokal menggunakan .env
	// Kode ini akan memeriksa apakah DATABASE_URL ada (lingkungan produksi)
	// Jika tidak, ia akan menggunakan variabel dari .env (lingkungan lokal)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		dsn = databaseURL
		log.Println("Connecting to production database...")
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
		log.Println("Connecting to local database...")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection successful.")

	err = DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	log.Println("Database migration successful.")
}
