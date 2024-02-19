package admin_test

import (
	"bytes"
	"encoding/json"
	Repo "my-us-stock-backend/app/repository/market-price/fund"
	"my-us-stock-backend/app/rest/admin"
	"my-us-stock-backend/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// テスト用のユーザーコントローラをセットアップ
func setupFundPriceController(db *gorm.DB) *admin.FundPriceController {
    fundPriceRepository := Repo.NewFetchFundRepository(db)
    fundPriceService := admin.NewFundPriceService(fundPriceRepository)
    return admin.NewFundPriceController(fundPriceService)
}

func TestGetFundPricesE2E(t *testing.T) {
    db := test.SetupTestDB() // Assume this function sets up your test database
    controller := setupFundPriceController(db) // Setup your controller here

    // Assuming you've inserted some test fund prices into the database here

    router := gin.Default()
    router.GET("/api/v1/admin/fund-prices", controller.GetFundPrices)

    req, _ := http.NewRequest("GET", "/api/v1/admin/fund-prices", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    // Further assertions to verify the response body can be added here
}

func TestCreateFundPriceE2E(t *testing.T) {
    db := test.SetupTestDB()
    controller := setupFundPriceController(db)

    router := gin.Default()
    router.POST("/api/v1/admin/fund-prices", controller.CreateFundPrice)

    newFundPrice := map[string]interface{}{
        "name": "New Fund",
        "code": "NF123",
        "price": 300.0,
    }
    body, _ := json.Marshal(newFundPrice)

    req, _ := http.NewRequest("POST", "/api/v1/admin/fund-prices", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    // Further assertions to verify the response body can be added here
}

func TestUpdateFundPriceE2E(t *testing.T) {
    db := test.SetupTestDB()
    controller := setupFundPriceController(db)

    // Assuming you've inserted a test fund price into the database here

    router := gin.Default()
    router.PUT("/api/v1/admin/fund-prices", controller.UpdateFundPrice)

    updatedFundPrice := map[string]interface{}{
        "price": 320.0,
    }
    body, _ := json.Marshal(updatedFundPrice)

    req, _ := http.NewRequest("PUT", "/api/v1/admin/fund-prices", bytes.NewBuffer(body)) // Assuming '1' is a valid ID
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    // Further assertions to verify the response body can be added here
}

