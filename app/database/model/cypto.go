package model

import (
	"my-us-stock-backend/app/repository/market-price/crypto"

	"gorm.io/gorm"
)

// Crypto は暗号通貨情報を表します。
type Crypto struct {
    gorm.Model
	Code   crypto.CryptoCode `gorm:"size:10;not null"`
	GetPrice float64 `gorm:"type:float"`
	Quantity float64 `gorm:"type:float"`
	UserId uint `gorm:"not null"`
}
