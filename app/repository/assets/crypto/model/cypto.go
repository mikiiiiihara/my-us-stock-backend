package model

import (
	"gorm.io/gorm"
)

// Crypto は暗号通貨情報を表します。
type Crypto struct {
    gorm.Model
	Ticker   string  `gorm:"size:10;not null"`
	GetPrice float64 `gorm:"type:float"`
	Quantity float64 `gorm:"type:float"`
}
