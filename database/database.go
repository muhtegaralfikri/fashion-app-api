package database


import (
    "fmt"
    "log"
    "appFashion/backend/models"
    "os"


    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)


// Variabel global untuk instance database
var DB *gorm.DB


// Fungsi untuk menghubungkan ke database
func ConnectDatabase() {
    var err error
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )


    // Membuka koneksi ke database
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database!")
    }


    log.Println("Database connection successful.")


    // AutoMigrate akan membuat tabel berdasarkan struct model Anda
    // Ini adalah cara cepat untuk sinkronisasi skema DB.
    err = DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{})
    if err != nil {
        log.Fatal("Failed to migrate database schema!")
    }
    log.Println("Database migration successful.")
}
