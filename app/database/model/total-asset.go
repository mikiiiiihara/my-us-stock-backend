package model

import (
	"gorm.io/gorm"
)

// TotalAsset は資産総額を表します。
type TotalAsset struct {
    gorm.Model
	CashJpy float64 `gorm:"type:float"`// 円ベースで登録
	CashUsd float64 `gorm:"type:float"` // ドルベースで登録
	Stock float64 `gorm:"type:float"`// 円ベースで登録
	Crypto float64 `gorm:"type:float"`// 円ベースで登録
	FixedIncomeAsset float64 `gorm:"type:float"`// 円ベースで登録
	UserId uint `gorm:"not null"`
}
