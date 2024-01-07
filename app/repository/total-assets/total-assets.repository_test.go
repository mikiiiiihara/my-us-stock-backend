package totalassets

import (
	"context"
	"errors"
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
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // UTCの日付の境界値でのテスト
    utcNow := time.Now().UTC()
    testBoundaryValues := []time.Time{
        utcNow.Add(-time.Hour),   // UTCの日付変更直前
        utcNow.Add(time.Hour),    // UTCの日付変更直後
        time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day(), 14, 59, 0, 0, time.UTC), // UTCで14時59分(JSTとUTCで日付が同一)
    }

    for _, testTime := range testBoundaryValues {
        asset := model.TotalAsset{
            Model: gorm.Model{
                CreatedAt: testTime,
            },
            UserId:  1,
            CashUsd: 1000,
        }
        db.Create(&asset)

        foundAsset, err := repo.FindTodayTotalAsset(context.Background(), asset.UserId)
        assert.NoError(t, err)
        assert.NotNil(t, foundAsset)
        assert.Equal(t, asset.UserId, foundAsset.UserId)
        assert.WithinDuration(t, testTime, foundAsset.CreatedAt, time.Second)

        // テストデータをクリーンアップ
        db.Unscoped().Delete(&asset)
    }

    // JSTとUTCとで日付が変わってしまう場合でのテスト
    pastAsset := model.TotalAsset{
        Model: gorm.Model{
            CreatedAt: utcNow.AddDate(0, 0, -1), // 昨日
        },
        UserId:  1,
        CashUsd: 1000,
    }
    futureAsset := model.TotalAsset{
        Model: gorm.Model{
            CreatedAt: utcNow.AddDate(0, 0, 1), // 明日
        },
        UserId:  1,
        CashUsd: 1000,
    }
    db.Create(&pastAsset)
    db.Create(&futureAsset)

    _, err := repo.FindTodayTotalAsset(context.Background(), pastAsset.UserId)
    assert.Error(t, err)
    _, err = repo.FindTodayTotalAsset(context.Background(), futureAsset.UserId)
    assert.Error(t, err)

    // DB初期化
    db.Unscoped().Where("1=1").Delete(&model.TotalAsset{})
}

func TestFindTodayTotalAssetNotFoundAtJST(t *testing.T) {
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // UTCで15時00分の時刻を設定((JSTとUTCで日付がズレる))
    utcNow := time.Now().UTC()
    utc1501 := time.Date(utcNow.Year(), utcNow.Month(), utcNow.Day(), 15, 0, 0, 0, time.UTC)

    // レコードを作成
    asset := model.TotalAsset{
        Model: gorm.Model{
            CreatedAt: utc1501,
        },
        UserId:  1,
        CashUsd: 1000,
    }
    db.Create(&asset)

    // 当日の資産総額を取得し、NotFoundエラーが発生することを確認
    _, err := repo.FindTodayTotalAsset(context.Background(), asset.UserId)
    assert.Error(t, err)

    // エラーの種類が `record not found` であることを確認
    assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

    // テストデータをクリーンアップ
    db.Unscoped().Delete(&asset)

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
