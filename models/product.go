package models


import "gorm.io/gorm"


type Category struct {
    gorm.Model
    Name string `gorm:"size:255;not null;unique" json:"name"`
}


type Product struct {
    gorm.Model
    Name          string  `gorm:"size:255;not null" json:"name"`
    Description   string  `json:"description"`
    Price         float64 `gorm:"not null" json:"price"`
    OriginalPrice float64 `json:"original_price"`
    ImageURL      string  `gorm:"size:255" json:"image_url"`
    Rating        float32 `json:"rating"`
    Reviews       int     `json:"reviews"`
    CategoryID    uint    `json:"category_id"`
    Category      Category `gorm:"foreignKey:CategoryID" json:"category"`
}


