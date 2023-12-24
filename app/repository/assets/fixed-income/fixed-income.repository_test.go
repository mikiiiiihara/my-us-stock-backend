package fixedincome

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/assets/fixed-income/dto"
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
    db.AutoMigrate(&model.FixedIncomeAsset{})

    return db
}

func TestFetchFixedIncomeAssetListById(t *testing.T) {
    db := setupTestDB()
    repo := NewFixedIncomeRepository(db)

    // テスト用データを作成
    fixedIncomeAsset := model.FixedIncomeAsset{Code: "Funds", UserId: "user1", DividendRate: 3.5, GetPriceTotal: 100000}
    db.Create(&fixedIncomeAsset)

    // User IDで検索
    fixedIncomeAssetList, err := repo.FetchFixedIncomeAssetListById(context.Background(), fixedIncomeAsset.UserId)
    assert.NoError(t, err)
    assert.NotEmpty(t, fixedIncomeAssetList)
    assert.Equal(t, fixedIncomeAsset.Code, fixedIncomeAssetList[0].Code)
    assert.Equal(t, fixedIncomeAsset.DividendRate, fixedIncomeAssetList[0].DividendRate)
	assert.Equal(t, fixedIncomeAsset.GetPriceTotal, fixedIncomeAssetList[0].GetPriceTotal)
	assert.Equal(t, fixedIncomeAsset.UserId, fixedIncomeAssetList[0].UserId)
}

// 取得結果が0件だった場合、空配列が返却される
func TestFetchFixedIncomeAssetListByIdEmpty(t *testing.T) {
    db := setupTestDB()
    repo := NewFixedIncomeRepository(db)

    // テスト用データを作成
    fixedIncomeAsset := model.FixedIncomeAsset{Code: "Funds", UserId: "user1", DividendRate: 3.5, GetPriceTotal: 100000}
    db.Create(&fixedIncomeAsset)

    // User IDで検索
    fixedIncomeAssetList, err := repo.FetchFixedIncomeAssetListById(context.Background(), "user2")
    assert.NoError(t, err)
    assert.Empty(t, fixedIncomeAssetList)
}

func TestUpdateFixedIncomeAsset(t *testing.T) {
    db := setupTestDB()
    repo := NewFixedIncomeRepository(db)

    // テスト用データを作成
    originalFixedIncomeAsset := model.FixedIncomeAsset{Code: "Funds", UserId: "user1", DividendRate: 3.5, GetPriceTotal: 100000}
    db.Create(&originalFixedIncomeAsset)

    // 更新用DTOの作成
    updateDto := dto.UpdateFixedIncomeDto{
        ID:       originalFixedIncomeAsset.ID,
        GetPriceTotal: new(float64),
    }
    *updateDto.GetPriceTotal = 115000

    // 株式情報を更新
    updatedFixedIncomeAsset, err := repo.UpdateFixedIncomeAsset(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedFixedIncomeAsset)
    assert.Equal(t, *updateDto.GetPriceTotal, updatedFixedIncomeAsset.GetPriceTotal)

	// 変わってないことを確認
	assert.Equal(t, originalFixedIncomeAsset.DividendRate, updatedFixedIncomeAsset.DividendRate)

    // データベースから直接取得して検証
    var dbFixedIncomeAsset model.FixedIncomeAsset
    db.First(&dbFixedIncomeAsset, originalFixedIncomeAsset.ID)
    assert.Equal(t, *updateDto.GetPriceTotal, dbFixedIncomeAsset.GetPriceTotal)
}

func TestCreateFixedIncomeAsset(t *testing.T) {
    db := setupTestDB()
    repo := NewFixedIncomeRepository(db)

    // 新しい株式情報を作成
    createDto := dto.CreateFixedIncomeDto{
        Code:   "Bankers",
        UserId:   "user1",
		DividendRate: 3.5, 
		GetPriceTotal: 100000,
    }
    created, err := repo.CreateFixedIncomeAsset(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, created)
    assert.Equal(t, createDto.Code, created.Code)
    assert.Equal(t, createDto.GetPriceTotal, created.GetPriceTotal)
	assert.Equal(t, createDto.DividendRate, created.DividendRate)
	assert.Equal(t, createDto.UsdJpy, created.UsdJpy)
	assert.Equal(t, createDto.UserId, created.UserId)

    // データベースで株式情報を確認
    var fixedIncomeAsset model.FixedIncomeAsset
    db.First(&fixedIncomeAsset, created.ID)
    assert.Equal(t, createDto.Code, fixedIncomeAsset.Code)
    assert.Equal(t, createDto.GetPriceTotal, fixedIncomeAsset.GetPriceTotal)
	assert.Equal(t, createDto.DividendRate, fixedIncomeAsset.DividendRate)
	assert.Equal(t, createDto.UsdJpy, fixedIncomeAsset.UsdJpy)
	assert.Equal(t, createDto.UserId, fixedIncomeAsset.UserId)
}

func TestCreateFixedIncomeAssetAlreadyExists(t *testing.T) {
    db := setupTestDB()
    repo := NewFixedIncomeRepository(db)

    // 既存の銘柄をデータベースに登録
    existingFixedIncomeAsset := model.FixedIncomeAsset{Code: "Funds", UserId: "user1", DividendRate: 3.5, GetPriceTotal: 100000}
    db.Create(&existingFixedIncomeAsset)

    // 同じ銘柄で新しい株式情報を作成
    createDto := dto.CreateFixedIncomeDto{
        Code:   "Funds",
        UserId:   "user1",
		DividendRate: 3.5, 
		GetPriceTotal: 100000,
    }

    // CreateUsStock メソッドを実行し、エラーが発生することを確認
    _, err := repo.CreateFixedIncomeAsset(context.Background(), createDto)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "この銘柄は既に登録されています")
}

func TestDeleteFixedIncomeAsset(t *testing.T) {
    db := setupTestDB()
    repo := NewFixedIncomeRepository(db)

    // テスト用データを作成
    fixedIncomeAsset  := model.FixedIncomeAsset{Code: "Funds", UserId: "user1", DividendRate: 3.5, GetPriceTotal: 100000}
    db.Create(&fixedIncomeAsset)

    // 株式情報を削除
    err := repo.DeleteFixedIncomeAsset(context.Background(), fixedIncomeAsset.ID)
    assert.NoError(t, err)

    // データベースから確認
    var result model.FixedIncomeAsset
    db.First(&result, fixedIncomeAsset.ID)
    assert.Empty(t, result)
}
