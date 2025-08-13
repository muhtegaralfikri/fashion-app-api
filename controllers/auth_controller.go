// File: controllers/auth_controller.go

package controllers

// NOTE: Jika Anda melihat error pada import di bawah (misalnya pada jwt),
// hentikan server (Ctrl+C) dan jalankan `go mod tidy` di terminal Anda
// untuk mengunduh paket yang dibutuhkan.
import (
	"net/http"
	"os"
	"time"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"

	"appFashion/backend/database" // Pastikan ini sesuai dengan nama modul Anda
	"appFashion/backend/models"    // Pastikan ini sesuai dengan nama modul Anda

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// DTO (Data Transfer Object) untuk input registrasi
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Fungsi untuk menangani request Registrasi
func Register(c *gin.Context) {
	var input RegisterInput

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Buat user baru
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Simpan user ke database
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}


// DTO untuk input Login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Fungsi untuk menangani request Login
func Login(c *gin.Context) {
	var input LoginInput

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Cari user berdasarkan email
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Bandingkan password yang diinput dengan hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Jika password cocok, generate token
	token, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Fungsi untuk membuat JWT Token
func GenerateToken(userID uint) (string, error) {
	// Menyiapkan klaim untuk token
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token berlaku selama 24 jam

	// Membuat token dengan klaim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key dari .env
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

// DTO untuk input Forgot Password
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

// Fungsi untuk menangani permintaan lupa password
func ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		// Kita tetap mengembalikan pesan sukses untuk tidak memberitahu penyerang
		// apakah sebuah email terdaftar atau tidak.
		c.JSON(http.StatusOK, gin.H{"message": "If a user with that email exists, a reset token has been generated."})
		return
	}

	// Generate token acak yang aman
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	resetToken := hex.EncodeToString(tokenBytes)

	// Set token dan waktu kedaluwarsa (misalnya 1 jam)
	user.ResetPasswordToken = resetToken
	user.ResetPasswordExpires = time.Now().Add(1 * time.Hour)

	database.DB.Save(&user)

	// PENTING: Di aplikasi nyata, Anda akan mengirim email di sini.
	// Untuk pengembangan, kita kembalikan tokennya agar bisa diuji.
	c.JSON(http.StatusOK, gin.H{
		"message": "Reset token generated. In a real app, this would be emailed.",
		"reset_token": resetToken,
	})
}


// DTO untuk input Reset Password
type ResetPasswordInput struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Fungsi untuk mereset password dengan token
func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Cari user berdasarkan token DAN pastikan token belum kedaluwarsa
	err := database.DB.Where("reset_password_token = ? AND reset_password_expires > ?", input.Token, time.Now()).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired reset token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update password dan hapus token agar tidak bisa dipakai lagi
	user.Password = string(hashedPassword)
	user.ResetPasswordToken = ""
	// Kita tidak perlu set ulang ResetPasswordExpires karena tokennya sudah tidak valid

	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully."})
}
