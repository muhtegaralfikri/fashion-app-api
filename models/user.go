package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                 string `gorm:"size:255;not null" json:"name"`
	Email                string `gorm:"size:255;not null;unique" json:"email"`
	Password             string `gorm:"size:255;not null" json:"-"`
	// --- TAMBAHKAN FIELD DI BAWAH INI ---
	ResetPasswordToken   string
	ResetPasswordExpires time.Time
}