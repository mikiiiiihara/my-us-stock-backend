package test

import (
	"my-us-stock-backend/app/database/model"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB はテスト用のデータベースをセットアップします
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.UsStock{})
	db.AutoMigrate(&model.Crypto{})
	db.AutoMigrate(&model.FixedIncomeAsset{})
	db.AutoMigrate(&model.JapanFund{})
	return db
}