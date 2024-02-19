package model

import (
	"gorm.io/gorm"
)

// FundPrice は投資信託の市場価格を表します。
type FundPrice struct {
    gorm.Model
	Name   string  `gorm:"size:255;not null"`
	Code   string  `gorm:"size:10;not null"`
	Price float64 `gorm:"type:float"`
}
