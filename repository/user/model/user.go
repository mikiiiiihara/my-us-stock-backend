package model

import (
    "gorm.io/gorm"
)

// User はユーザー情報を表します。
type User struct {
    gorm.Model
    Name  string
    Email string `gorm:"uniqueIndex"`
}
