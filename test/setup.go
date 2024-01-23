package test

import (
	"log"
	"my-us-stock-backend/app/database/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB はテスト用のデータベースをセットアップします
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.UsStock{})
	db.AutoMigrate(&model.Crypto{})
	db.AutoMigrate(&model.FixedIncomeAsset{})
	db.AutoMigrate(&model.JapanFund{})
	db.AutoMigrate(&model.TotalAsset{})
	return db
}