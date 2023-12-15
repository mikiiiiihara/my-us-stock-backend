package database

import (
	"log"
	CryptoModel "my-us-stock-backend/app/repository/assets/crypto/model"
	FixedIncomeAssetModel "my-us-stock-backend/app/repository/assets/fixed-income/model"
	JapanFundModel "my-us-stock-backend/app/repository/assets/fund/japan/model"
	UsStockModel "my-us-stock-backend/app/repository/assets/stock/us/model"
	StrategyModel "my-us-stock-backend/app/repository/strategy/model"
	UserModel "my-us-stock-backend/app/repository/user/model"
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
	db.AutoMigrate(&CryptoModel.Crypto{})
	db.AutoMigrate(&FixedIncomeAssetModel.FixedIncomeAsset{})
	db.AutoMigrate(&JapanFundModel.JapanFund{})
	db.AutoMigrate(&UsStockModel.UsStock{})
	db.AutoMigrate(&StrategyModel.Strategy{})
	db.AutoMigrate(&UserModel.User{})
}
