package model

import (
	"gorm.io/gorm"
)

// JapanFund は日本の投資信託(三菱UFJのみ)を表します。
type JapanFund struct {
    gorm.Model
	Name   string  `gorm:"size:255;not null"`
	Code   string  `gorm:"size:10;not null"`
	GetPriceTotal float64 `gorm:"type:float"`
	GetPrice float64 `gorm:"type:float"`
	UserId uint `gorm:"not null;index"`
}
