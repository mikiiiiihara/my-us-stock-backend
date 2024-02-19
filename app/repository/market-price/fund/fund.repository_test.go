package fund

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB similar to the one provided but ensures model.FundPrice is migrated.
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    db.AutoMigrate(&model.FundPrice{})
    return db
}

// TestFetchFundPriceList tests FetchFundPriceList function of DefaultFundPriceRepository.
func TestFetchFundPriceList(t *testing.T) {
    db := setupTestDB()
    repo := NewFetchFundRepository(db)

    // Insert test data
    db.Create(&model.FundPrice{Name: "Test Fund", Code: "TF123", Price: 100.0})

    // Test FetchFundPriceList
    prices, err := repo.FetchFundPriceList(context.Background())
    assert.NoError(t, err)
    assert.NotEmpty(t, prices)
    // Add more assertions as needed
}
// TestFindFundPriceByCode tests the FindFundPriceByCode function of DefaultFundPriceRepository.
func TestFindFundPriceByCode(t *testing.T) {
    db := setupTestDB()
    repo := NewFetchFundRepository(db)

    // Setup test data
    expectedFundPrice := model.FundPrice{Name: "Test Fund", Code: "TF123", Price: 100.0}
    db.Create(&expectedFundPrice)

    // Test FindFundPriceByCode with an existing code
    foundFundPrice, err := repo.FindFundPriceByCode(context.Background(), "TF123")
    assert.NoError(t, err)
    assert.NotNil(t, foundFundPrice)
    assert.Equal(t, expectedFundPrice.Name, foundFundPrice.Name)
    assert.Equal(t, expectedFundPrice.Code, foundFundPrice.Code)
    assert.Equal(t, expectedFundPrice.Price, foundFundPrice.Price)

    // Test FindFundPriceByCode with a non-existing code
    _, err = repo.FindFundPriceByCode(context.Background(), "NON_EXISTENT_CODE")
    assert.Error(t, err)
    assert.Equal(t, gorm.ErrRecordNotFound, err)
}

// TestUpdateFundPrice tests UpdateFundPrice function of DefaultFundPriceRepository.
func TestUpdateFundPrice(t *testing.T) {
    db := setupTestDB()
    repo := NewFetchFundRepository(db)

    // Setup and insert test data
    originalPrice := model.FundPrice{Name: "Test Fund", Code: "TF123", Price: 100.0}
    db.Create(&originalPrice)

    // Update test data
    updateDto := UpdateFundPriceDto{ID: originalPrice.ID, Price: 105.0}
    updatedPrice, err := repo.UpdateFundPrice(context.Background(), updateDto)
    assert.NoError(t, err)
    assert.NotNil(t, updatedPrice)
    assert.Equal(t, 105.0, updatedPrice.Price)
    // Add more assertions as needed
}

// TestCreateFundPrice tests CreateFundPrice function of DefaultFundPriceRepository.
func TestCreateFundPrice(t *testing.T) {
    db := setupTestDB()
    repo := NewFetchFundRepository(db)

    // Create a new fund price
    createDto := CreateFundPriceDto{Name: "New Fund", Code: "NF123", Price: 110.0}
    createdPrice, err := repo.CreateFundPrice(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, createdPrice)
    // Add more assertions as needed

    // Test for error when creating a fund price with an existing code
    _, err = repo.CreateFundPrice(context.Background(), createDto)
    assert.Error(t, err)
}

func TestCreateFundPriceAlreadyExists(t *testing.T) {
    db := setupTestDB()
    repo := NewFetchFundRepository(db)

	// 既存の銘柄をデータベースに登録
	db.Create(&model.FundPrice{Name: "Test Fund", Code: "TF1234", Price: 100.0})
    // Create a new fund price
    createDto := CreateFundPriceDto{Name: "Test Fund", Code: "TF1234", Price: 110.0}
    _, err := repo.CreateFundPrice(context.Background(), createDto)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "この銘柄は既に登録されています")
}
