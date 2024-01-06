package totalassets

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"testing"
	"time"

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
    db.AutoMigrate(&model.TotalAsset{})

    return db
}

// FetchTotalAssetListByIdのテスト
func TestFetchTotalAssetListById(t *testing.T) {
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // テストデータの作成
    asset := model.TotalAsset{UserId: 1, CashUsd: 1000}
    db.Create(&asset)

    // 正常に取得できることを確認
    assets, err := repo.FetchTotalAssetListById(context.Background(), asset.UserId, 7)
    assert.NoError(t, err)
    assert.NotEmpty(t, assets)
    assert.Equal(t, asset.UserId, assets[0].UserId)
    assert.Equal(t, asset.CashUsd, assets[0].CashUsd)

    // 存在しないユーザーIDで検索
    emptyAssets, err := repo.FetchTotalAssetListById(context.Background(), 999, 7)
    assert.NoError(t, err)
    assert.Empty(t, emptyAssets)
	// DB初期化
	db.Unscoped().Where("1=1").Delete(&model.TotalAsset{})
}

func TestFindTodayTotalAsset(t *testing.T) {
    // テスト実行前にタイムゾーンをUTCに設定
    time.Local = time.UTC
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // テストデータの作成
    currentUTC := time.Now().UTC()
    asset := model.TotalAsset{
        Model: gorm.Model{
            CreatedAt: currentUTC,
        },
        UserId:  1,
        CashUsd: 1000,
    }
    db.Create(&asset)

    // 当日の資産総額を取得
    foundAsset, err := repo.FindTodayTotalAsset(context.Background(), asset.UserId)
    assert.NoError(t, err)
    assert.NotNil(t, foundAsset)
    assert.Equal(t, asset.UserId, foundAsset.UserId)
    assert.Equal(t, currentUTC, foundAsset.CreatedAt)

    // 存在しないユーザーIDで検索
    _, err = repo.FindTodayTotalAsset(context.Background(), 999)
    assert.Error(t, err)

    // DB初期化
    db.Unscoped().Where("1=1").Delete(&model.TotalAsset{})
}

// UpdateTotalAssetのテスト
func TestUpdateTotalAsset(t *testing.T) {
    // テスト実行前にタイムゾーンをUTCに設定
    time.Local = time.UTC
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // テストデータの作成
    asset := model.TotalAsset{UserId: 1, CashUsd: 1000}
    db.Create(&asset)

    // 更新用DTOの作成
    updateDto := UpdateTotalAssetDto{ID: asset.ID, CashUsd: new(float64)}
    *updateDto.CashUsd = 1500

    // 資産情報を更新
    updatedAsset, err := repo.UpdateTotalAsset(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedAsset)
    assert.Equal(t, *updateDto.CashUsd, updatedAsset.CashUsd)

    // 存在しないIDで更新
    invalidUpdateDto := UpdateTotalAssetDto{ID: 999, CashUsd: new(float64)}
    *invalidUpdateDto.CashUsd = 2000
    _, err = repo.UpdateTotalAsset(context.Background(), invalidUpdateDto)
    assert.Error(t, err)
	// DB初期化
	db.Unscoped().Where("1=1").Delete(&model.TotalAsset{})
}

// CreateTodayTotalAssetのテスト
func TestCreateTodayTotalAsset(t *testing.T) {
    // テスト実行前にタイムゾーンをUTCに設定
    time.Local = time.UTC
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // 新規資産の作成
    createDto := CreateTotalAssetDto{UserId: 1, CashUsd: 1000}
    createdAsset, err := repo.CreateTodayTotalAsset(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, createdAsset)
    // 同じ日付での重複作成
    _, err = repo.CreateTodayTotalAsset(context.Background(), createDto)
    assert.Error(t, err)
	// エラーメッセージをチェックする
	assert.Equal(t, "既に資産が登録されています。新規追加ではなく更新を行ってください。", err.Error())
	// DB初期化
	db.Unscoped().Where("1=1").Delete(&model.TotalAsset{})
}
