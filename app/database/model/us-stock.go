package model

import (
	"gorm.io/gorm"
)

// UsStock は米国株式を表します。
type UsStock struct {
    gorm.Model
	Code   string  `gorm:"size:6;not null"`
	GetPrice float64 `gorm:"type:float"`
	Quantity float64 `gorm:"type:float"`
	Sector   string  `gorm:"size:255;not null"`
	UsdJpy   float64 `gorm:"type:float"`
	UserId uint `gorm:"not null"`
}
