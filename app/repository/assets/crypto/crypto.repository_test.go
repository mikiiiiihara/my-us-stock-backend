package crypto

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/assets/crypto/dto"
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
    db.AutoMigrate(&model.Crypto{})

    return db
}

func TestFetchCryptoListById(t *testing.T) {
    db := setupTestDB()
    repo := NewCryptoRepository(db)

    // テスト用データを作成
    crypto := model.Crypto{Code: "xrp", UserId: 99, Quantity: 100, GetPrice: 80}
    db.Create(&crypto)

    // User IDで検索
    cryptoList, err := repo.FetchCryptoListById(context.Background(), crypto.UserId)
    assert.NoError(t, err)
    assert.NotEmpty(t, cryptoList)
    assert.Equal(t, crypto.Code, cryptoList[0].Code)
    assert.Equal(t, crypto.Quantity, cryptoList[0].Quantity)
}

// 取得結果が0件だった場合、空配列が返却される
func TestFetchCryptoListByIdEmpty(t *testing.T) {
    db := setupTestDB()
    repo := NewCryptoRepository(db)

    // テスト用データを作成
    crypto := model.Crypto{Code: "xrp", UserId: 99, Quantity: 100, GetPrice: 80}
    db.Create(&crypto)

    // User IDで検索
    cryptoList, err := repo.FetchCryptoListById(context.Background(), 98)
    assert.NoError(t, err)
    assert.Empty(t, cryptoList)
}


func TestUpdateCrypto(t *testing.T) {
    db := setupTestDB()
    repo := NewCryptoRepository(db)

    // テスト用データを作成
    originalCrypto :=model.Crypto{Code: "xrp", UserId: 99, Quantity: 100, GetPrice: 80}
    db.Create(&originalCrypto)

    // 更新用DTOの作成
    updateDto := dto.UpdateCryptoDto{
        ID:       originalCrypto.ID,
        Quantity: new(float64),
    }
    *updateDto.Quantity = 115

    // 株式情報を更新
    updatedCrypto, err := repo.UpdateCrypto(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedCrypto)
    assert.Equal(t, *updateDto.Quantity, updatedCrypto.Quantity)

		// 変わってないことを確認
		assert.Equal(t, originalCrypto.GetPrice, updatedCrypto.GetPrice)

    // データベースから直接取得して検証
    var dbCrypto model.Crypto
    db.First(&dbCrypto, originalCrypto.ID)
    assert.Equal(t, *updateDto.Quantity, dbCrypto.Quantity)
}

func TestCryptoUsStock(t *testing.T) {
    db := setupTestDB()
    repo := NewCryptoRepository(db)

    // 新しい株式情報を作成
    createDto := dto.CreateCryptDto{
        Code:   "btc",
        UserId:   99,
        Quantity: 0.4,
    }
    created, err := repo.CreateCrypto(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, created)
    assert.Equal(t, createDto.Code, created.Code)
    assert.Equal(t, createDto.Quantity, created.Quantity)

    // データベースで株式情報を確認
    var crypto model.Crypto
    db.First(&crypto, created.ID)
    assert.Equal(t, createDto.Code, crypto.Code)
    assert.Equal(t, createDto.Quantity, crypto.Quantity)
}

func TestCreateCryptoAlreadyExists(t *testing.T) {
    db := setupTestDB()
    repo := NewCryptoRepository(db)

    // 既存の銘柄をデータベースに登録
    existingCrypto := model.Crypto{Code: "xrp", UserId: 99, Quantity: 100, GetPrice: 80}
    db.Create(&existingCrypto)

    // 同じ銘柄で新しい株式情報を作成
    createDto := dto.CreateCryptDto{
        Code:   "xrp",
        UserId:   99,
        Quantity: 20,
        GetPrice: 86.0,
    }

    // CreateUsStock メソッドを実行し、エラーが発生することを確認
    _, err := repo.CreateCrypto(context.Background(), createDto)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "この銘柄は既に登録されています")
}

func TestDeleteUsStock(t *testing.T) {
    db := setupTestDB()
    repo := NewCryptoRepository(db)

    // テスト用データを作成
    crypto  := model.Crypto{Code: "xrp", UserId: 99, Quantity: 100, GetPrice: 80}
    db.Create(&crypto)

    // 株式情報を削除
    err := repo.DeleteCrypto(context.Background(), crypto.ID)
    assert.NoError(t, err)

    // データベースから確認
    var result model.Crypto
    db.First(&result, crypto.ID)
    assert.Empty(t, result)
}
