// File: controllers/product_controller.go

package controllers

import (
	"net/http"

	"appFashion/backend/database"
	"appFashion/backend/models"

	"github.com/gin-gonic/gin"
)

// Fungsi untuk mendapatkan semua produk
// Mendukung filter berdasarkan kategori
func GetProducts(c *gin.Context) {
	var products []models.Product
	
	// Cek apakah ada query parameter "category_id"
	categoryID := c.Query("category_id")

	// Buat query dasar
	query := database.DB.Model(&models.Product{})

	// Jika ada filter kategori, tambahkan ke query
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// Eksekusi query
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Fungsi untuk mendapatkan satu produk berdasarkan ID
func GetProductByID(c *gin.Context) {
	var product models.Product

	// Ambil ID dari URL parameter
	id := c.Param("id")

	// Cari produk di database
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
