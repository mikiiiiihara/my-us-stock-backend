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
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // UTCの日付の境界値でのテスト
    utcNow := time.Now().UTC()
    testBoundaryValues := []time.Time{
        utcNow.Add(-time.Hour),   // UTCの日付変更直前
        utcNow.Add(time.Hour),    // UTCの日付変更直後
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

func TestFindTodayTotalAssetInJST(t *testing.T) {
    // テスト実行前にタイムゾーンをUTCに設定
    time.Local = time.UTC
    db := setupTestDB()
    repo := NewTotalAssetRepository(db)

    // JSTタイムゾーンの作成 (+9時間)
    jst := time.FixedZone("JST", 9*60*60)

    // JSTでの現在時刻
    jstNow := time.Now().In(jst)

    // JSTでのテストデータ作成
    jstAsset := model.TotalAsset{
        Model: gorm.Model{
            CreatedAt: jstNow, // JSTでの現在時刻
        },
        UserId:  1,
        CashUsd: 1000,
    }
    db.Create(&jstAsset)

    // UTC基準でデータ取得を試みる
    foundAsset, err := repo.FindTodayTotalAsset(context.Background(), jstAsset.UserId)
    assert.NoError(t, err)
    assert.NotNil(t, foundAsset)
    assert.Equal(t, jstAsset.UserId, foundAsset.UserId)

    // JSTのデータがUTC日付として正しく扱われているかを確認
    jstDay := jstNow.Format("2006-01-02")
    utcDay := foundAsset.CreatedAt.UTC().Format("2006-01-02")
    assert.Equal(t, jstDay, utcDay)

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
