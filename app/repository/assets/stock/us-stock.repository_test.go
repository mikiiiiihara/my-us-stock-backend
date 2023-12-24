package stock

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/assets/stock/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テスト用のデータベース設定
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // テスト用のテーブルを準備
    db.AutoMigrate(&model.UsStock{})

    return db
}

func TestFetchUsStockListById(t *testing.T) {
    db := setupTestDB()
    repo := NewUsStockRepository(db)

    // テスト用データを作成
    stock := model.UsStock{Code: "AAPL", UserId: "user1", Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9}
    db.Create(&stock)

    // User IDで検索
    stocks, err := repo.FetchUsStockListById(context.Background(), stock.UserId)
    assert.NoError(t, err)
    assert.NotEmpty(t, stocks)
    assert.Equal(t, stock.Code, stocks[0].Code)
    assert.Equal(t, stock.Quantity, stocks[0].Quantity)
}

// 取得結果が0件だった場合、空配列が返却される
func TestFetchUsStockListByIdEmpty(t *testing.T) {
    db := setupTestDB()
    repo := NewUsStockRepository(db)

    // テスト用データを作成
    stock := model.UsStock{Code: "AAPL", UserId: "user1", Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9}
    db.Create(&stock)

    // User IDで検索
    stocks, err := repo.FetchUsStockListById(context.Background(), "user2")
    assert.NoError(t, err)
    assert.Empty(t, stocks)
}


func TestUpdateUsStock(t *testing.T) {
    db := setupTestDB()
    repo := NewUsStockRepository(db)

    // テスト用データを作成
    originalStock := model.UsStock{Code: "AAPL", UserId: "user1", Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9}
    db.Create(&originalStock)

    // 更新用DTOの作成
    updateDto := dto.UpdateUsStockDto{
        ID:       originalStock.ID,
        Quantity: new(float64),
    }
    *updateDto.Quantity = 15.0

    // 株式情報を更新
    updatedStock, err := repo.UpdateUsStock(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedStock)
    assert.Equal(t, *updateDto.Quantity, updatedStock.Quantity)

    // データベースから直接取得して検証
    var dbStock model.UsStock
    db.First(&dbStock, originalStock.ID)
    assert.Equal(t, *updateDto.Quantity, dbStock.Quantity)
}

func TestCreateUsStock(t *testing.T) {
    db := setupTestDB()
    repo := NewUsStockRepository(db)

    // 新しい株式情報を作成
    createDto := dto.CreateUsStockDto{
        Code:   "MSFT",
        UserId:   "user1",
        Quantity: 5.0,
    }
    created, err := repo.CreateUsStock(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, created)
    assert.Equal(t, createDto.Code, created.Code)
    assert.Equal(t, createDto.Quantity, created.Quantity)

    // データベースで株式情報を確認
    var stock model.UsStock
    db.First(&stock, created.ID)
    assert.Equal(t, createDto.Code, stock.Code)
    assert.Equal(t, createDto.Quantity, stock.Quantity)
}

func TestCreateUsStockAlreadyExists(t *testing.T) {
    db := setupTestDB()
    repo := NewUsStockRepository(db)

    // 既存の銘柄をデータベースに登録
    existingStock := model.UsStock{Code: "AAPL", UserId: "user1", Quantity: 10, GetPrice: 150.0, Sector: "Tech", UsdJpy: 110.0}
    db.Create(&existingStock)

    // 同じ銘柄で新しい株式情報を作成
    createDto := dto.CreateUsStockDto{
        Code:   "AAPL",
        UserId:   "user1",
        Quantity: 20,
        GetPrice: 155.0,
        Sector:   "Tech",
        UsdJpy:   111.0,
    }

    // CreateUsStock メソッドを実行し、エラーが発生することを確認
    _, err := repo.CreateUsStock(context.Background(), createDto)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "この銘柄は既に登録されています")
}

func TestDeleteUsStock(t *testing.T) {
    db := setupTestDB()
    repo := NewUsStockRepository(db)

    // テスト用データを作成
    stock := model.UsStock{Code: "AAPL", UserId: "user1", Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9}
    db.Create(&stock)

    // 株式情報を削除
    err := repo.DeleteUsStock(context.Background(), stock.ID)
    assert.NoError(t, err)

    // データベースから確認
    var result model.UsStock
    db.First(&result, stock.ID)
    assert.Empty(t, result)
}
