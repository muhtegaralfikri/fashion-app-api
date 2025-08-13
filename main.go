package main

import (
	"log"
	"os"

	"appFashion/backend/database" // Pastikan ini sesuai dengan nama modul Anda
	"appFashion/backend/routes"    // Pastikan ini sesuai dengan nama modul Anda

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// PERBAIKAN: Membuat pemuatan .env menjadi opsional.
	// Ini hanya akan berjalan di lingkungan lokal.
	// Di Render, ini akan gagal secara diam-diam dan aplikasi akan tetap berjalan.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file. Using environment variables from host.")
	}

	// Inisialisasi koneksi database
	database.ConnectDatabase()

	// Membuat router Gin
	r := gin.Default()

	// Mengatur rute dari file routes/routes.go
	routes.SetupRoutes(r)

	// Menjalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Port default jika tidak diset di .env
	}
	r.Run(":" + port)
}
