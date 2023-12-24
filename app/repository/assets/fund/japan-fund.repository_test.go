package fund

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/assets/fund/dto"
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
    db.AutoMigrate(&model.JapanFund{})

    return db
}

func TestFetchJapanFundListById(t *testing.T) {
    db := setupTestDB()
    repo := NewJapanFundRepository(db)

    // テスト用データを作成
    fund := model.JapanFund{Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", UserId: "user1", Code: "253266", GetPrice: 15523.81, GetPriceTotal: 761157}
    db.Create(&fund)

    // User IDで検索
    funds, err := repo.FetchJapanFundListById(context.Background(), fund.UserId)
    assert.NoError(t, err)
    assert.NotEmpty(t, funds)
    assert.Equal(t, fund.Code, funds[0].Code)
    assert.Equal(t, fund.Name, funds[0].Name)
	assert.Equal(t, fund.GetPriceTotal, funds[0].GetPriceTotal)
	assert.Equal(t, fund.GetPrice, funds[0].GetPrice)
	assert.Equal(t, fund.UserId, funds[0].UserId)
}

// 取得結果が0件だった場合、空配列が返却される
func TestFetchJapanFundListByIdEmpty(t *testing.T) {
    db := setupTestDB()
    repo := NewJapanFundRepository(db)

    // テスト用データを作成
    fund := model.JapanFund{Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", UserId: "user1", Code: "253266", GetPrice: 15523.81, GetPriceTotal: 761157}
    db.Create(&fund)

    // User IDで検索
    funds, err := repo.FetchJapanFundListById(context.Background(), "user2")
    assert.NoError(t, err)
    assert.Empty(t, funds)
}

func TestUpdateJapanFund(t *testing.T) {
    db := setupTestDB()
    repo := NewJapanFundRepository(db)

    // テスト用データを作成
    originalFund := model.JapanFund{Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", UserId: "user1", Code: "253266", GetPrice: 15523.81, GetPriceTotal: 761157}
    db.Create(&originalFund)

    // 更新用DTOの作成
    updateDto := dto.UpdateJapanFundDto{
        ID:       originalFund.ID,
        GetPrice: new(float64),
    }
    *updateDto.GetPrice = 16000.0

    // 株式情報を更新
    updatedFund, err := repo.UpdateJapanFund(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedFund)
    assert.Equal(t, *updateDto.GetPrice, updatedFund.GetPrice)
	// 変わってないことを確認
	assert.Equal(t, originalFund.GetPriceTotal, updatedFund.GetPriceTotal)

    // データベースから直接取得して検証
    var dbFund model.JapanFund
    db.First(&dbFund, originalFund.ID)
    assert.Equal(t, *updateDto.GetPrice, dbFund.GetPrice)
}

func TestCreateJapanFund(t *testing.T) {
    db := setupTestDB()
    repo := NewJapanFundRepository(db)

    // 新しい株式情報を作成
    createDto := dto.CreateJapanFundDto{
        Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）",
        UserId: "user1",
        Code: "253267",
        GetPrice: 15523.81,
        GetPriceTotal: 761157,
    }
    created, err := repo.CreateJapanFund(context.Background(), createDto)

    // エラーが発生していないことを確認
    assert.NoError(t, err)

    // created が nil でないことを確認
    assert.NotNil(t, created)

    // その他のアサーションを実行
    assert.Equal(t, createDto.Name, created.Name)
    assert.Equal(t, createDto.Code, created.Code)
    assert.Equal(t, createDto.GetPriceTotal, created.GetPriceTotal)
    assert.Equal(t, createDto.GetPrice, created.GetPrice)
    assert.Equal(t, createDto.UserId, created.UserId)

    // データベースで株式情報を確認
    var fund model.JapanFund
    db.First(&fund, created.ID)
    assert.Equal(t, createDto.Name, fund.Name)
    assert.Equal(t, createDto.Code, fund.Code)
    assert.Equal(t, createDto.GetPriceTotal, fund.GetPriceTotal)
    assert.Equal(t, createDto.GetPrice, fund.GetPrice)
    assert.Equal(t, createDto.UserId, fund.UserId)
}

func TestCreateJapanFundAlreadyExists(t *testing.T) {
    db := setupTestDB()
    repo := NewJapanFundRepository(db)

    // 既存の銘柄をデータベースに登録
    existingFund := model.JapanFund{Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", UserId: "user1", Code: "253266", GetPrice: 15523.81, GetPriceTotal: 761157}
    db.Create(&existingFund)

    // 同じ銘柄で新しい情報を作成
    createDto := dto.CreateJapanFundDto{
		Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）",
		UserId: "user1",
		Code: "253266",
		GetPrice: 15523.81,
		GetPriceTotal: 761157,
    }

    // Createメソッドを実行し、エラーが発生することを確認
    _, err := repo.CreateJapanFund(context.Background(), createDto)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "この銘柄は既に登録されています")
}

func TestDeleteJapanFund(t *testing.T) {
    db := setupTestDB()
    repo := NewJapanFundRepository(db)

    // テスト用データを作成
    fund := model.JapanFund{Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", UserId: "user1", Code: "253266", GetPrice: 15523.81, GetPriceTotal: 761157}
    db.Create(&fund)

    // 株式情報を削除
    err := repo.DeleteJapanFund(context.Background(), fund.ID)
    assert.NoError(t, err)

    // データベースから確認
    var result model.JapanFund
    db.First(&result, fund.ID)
    assert.Empty(t, result)
}
