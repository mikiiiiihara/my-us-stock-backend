package model

import (
	"time"

	"gorm.io/gorm"
)

// FixedIncomeFund は債券・不動産・クラウドファンディングなど、固定された収入や利回りを持つ資産を表します。
type FixedIncomeAsset struct {
    gorm.Model
	Code   string  `gorm:"size:255;not null"`
	GetPriceTotal float64 `gorm:"type:float"`
	DividendRate float64 `gorm:"type:float"`
	UsdJpy   float64 `gorm:"type:float"`
	PaymentDate time.Time  `gorm:"type:timestamp"`
	UserId uint `gorm:"not null"`
}
