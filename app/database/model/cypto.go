package model

import (
	"gorm.io/gorm"
)

// Crypto は暗号通貨情報を表します。
type Crypto struct {
    gorm.Model
	Code   string  `gorm:"size:10;not null"`
	GetPrice float64 `gorm:"type:float"`
	Quantity float64 `gorm:"type:float"`
	UserId string `gorm:"size:255;not null"`
}
