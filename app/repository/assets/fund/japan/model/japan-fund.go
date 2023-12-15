package model

import (
	"gorm.io/gorm"
)

// JapanFund は日本の投資信託(三菱UFJのみ)を表します。
type JapanFund struct {
    gorm.Model
	Code   string  `gorm:"size:10;not null"`
	GetPriceTotal float64 `gorm:"type:float"`
}
