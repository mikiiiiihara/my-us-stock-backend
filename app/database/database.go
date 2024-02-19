package database

import (
	"log"
	"my-us-stock-backend/app/database/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect はデータベース接続を確立し、gorm.DBインスタンスを返します。
func Connect() *gorm.DB {
	dsn :=os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

// モデルに基づいてテーブルを作成または更新
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Crypto{})
	db.AutoMigrate(&model.TotalAsset{})
	db.AutoMigrate(&model.FixedIncomeAsset{})
	db.AutoMigrate(&model.JapanFund{})
	db.AutoMigrate(&model.UsStock{})
	db.AutoMigrate(&model.FundPrice{})
	db.AutoMigrate(&model.User{})
}
