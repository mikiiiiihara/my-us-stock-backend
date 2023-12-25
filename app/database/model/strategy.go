package model

import (
	"gorm.io/gorm"
)

// Strategy はユーザーが利用できる戦略メモを表します。
type Strategy struct {
    gorm.Model
    Text  string `gorm:"size:255;not null"`
	UserId uint `gorm:"not null"`
}
