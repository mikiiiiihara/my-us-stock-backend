package stock

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestUsStocks は UsStocks メソッドのテストです。
func TestUsStocksService(t *testing.T) {
	mockStockRepo := stock.NewMockUsStockRepository()
    mockMarketPriceRepo := marketPrice.NewMockMarketPriceRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewUsStockService(mockStockRepo, mockAuth, mockMarketPriceRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockStocks := []model.UsStock{
		{Code: "AAPL", GetPrice: 150, Quantity: 10, Sector: "IT"},
	}
	mockStockRepo.On("FetchUsStockListById", mock.Anything, userId).Return(mockStocks, nil)

	mockMarketPrices := []marketPrice.MarketPriceDto{
		{Ticker: "AAPL", CurrentPrice: 155, PriceGets: 5, CurrentRate: 0.0333},
	}
	mockMarketPriceRepo.On("FetchMarketPriceList", mock.Anything, []string{"AAPL"}).Return(mockMarketPrices, nil)

	mockDividend := &marketPrice.DividendEntity{DividendTotal: 1.5}
	mockMarketPriceRepo.On("FetchDividend", mock.Anything, "AAPL").Return(mockDividend, nil)

	// テスト対象メソッドの実行
	usStocks, err := service.UsStocks(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, usStocks)
	assert.Len(t, usStocks, 1)

	assert.Equal(t, "0", usStocks[0].ID)
	assert.Equal(t, "AAPL", usStocks[0].Code)
	assert.Equal(t, 150.0, usStocks[0].GetPrice)
	assert.Equal(t, 10.0, usStocks[0].Quantity)
	assert.Equal(t, "IT", usStocks[0].Sector)
	assert.Equal(t, 1.5, usStocks[0].Dividend)
	assert.Equal(t, 155.0, usStocks[0].CurrentPrice)
	assert.Equal(t, 5.0, usStocks[0].PriceGets)
	assert.Equal(t, 0.0333, usStocks[0].CurrentRate)

	// モックの呼び出しを検証
	mockStockRepo.AssertExpectations(t)
	mockMarketPriceRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestCreateUsStockService は TestCreateUsStock メソッドのテストです。
func TestCreateUsStockService(t *testing.T) {
	mockStockRepo := stock.NewMockUsStockRepository()
    mockMarketPriceRepo := marketPrice.NewMockMarketPriceRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewUsStockService(mockStockRepo, mockAuth, mockMarketPriceRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockStock := &model.UsStock{Code: "AAPL", GetPrice: 150, Quantity: 10, Sector: "IT",UsdJpy: 133.0}

	input := stock.CreateUsStockDto{
		Code: "AAPL",
		GetPrice: 150, 
		Quantity: 10, 
		Sector: "IT",
		UsdJpy: 133.0,
		UserId: 1,
	}
	mockStockRepo.On("CreateUsStock", mock.Anything, input).Return(mockStock, nil)

	mockMarketPrices := []marketPrice.MarketPriceDto{
		{Ticker: "AAPL", CurrentPrice: 155, PriceGets: 5, CurrentRate: 0.0333},
	}
	mockMarketPriceRepo.On("FetchMarketPriceList", mock.Anything, []string{"AAPL"}).Return(mockMarketPrices, nil)

	mockDividend := &marketPrice.DividendEntity{DividendTotal: 1.5}
	mockMarketPriceRepo.On("FetchDividend", mock.Anything, "AAPL").Return(mockDividend, nil)

	// テスト対象メソッドの実行
	serviceInput := generated.CreateUsStockInput{
		Code: "AAPL",
		GetPrice: 150, 
		Quantity: 10, 
		Sector: "IT",
		UsdJpy: 133.0,
	}
	usStock, err := service.CreateUsStock(context.Background(), serviceInput)
	assert.NoError(t, err)
	assert.NotNil(t, usStock)

	assert.Equal(t, "0", usStock.ID)
	assert.Equal(t, "AAPL", usStock.Code)
	assert.Equal(t, 150.0, usStock.GetPrice)
	assert.Equal(t, 10.0, usStock.Quantity)
	assert.Equal(t, "IT", usStock.Sector)
	assert.Equal(t, 1.5, usStock.Dividend)
	assert.Equal(t, 155.0, usStock.CurrentPrice)
	assert.Equal(t, 5.0, usStock.PriceGets)
	assert.Equal(t, 0.0333, usStock.CurrentRate)

	// モックの呼び出しを検証
	mockStockRepo.AssertExpectations(t)
	mockMarketPriceRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestUpdateUsStockService は UpdateUsStock メソッドのテストです。
func TestUpdateUsStockService(t *testing.T) {
	mockStockRepo := stock.NewMockUsStockRepository()
	mockMarketPriceRepo := marketPrice.NewMockMarketPriceRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewUsStockService(mockStockRepo, mockAuth, mockMarketPriceRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	updateInput := stock.UpdateUsStockDto{
		ID:        1,
		GetPrice:  new(float64),
		Quantity:  new(float64),
		UsdJpy:    new(float64),
	}
	*updateInput.GetPrice = 160
	*updateInput.Quantity = 20
	*updateInput.UsdJpy = 135.0

	updatedMockStock := &model.UsStock{
		Model: gorm.Model{
			ID: 1,
		},
		Code:      "AAPL",
		GetPrice:  160,
		Quantity:  20,
		Sector:    "IT",
		UsdJpy:    135.0,
	}
	mockStockRepo.On("UpdateUsStock", mock.Anything, updateInput).Return(updatedMockStock, nil)

	mockMarketPrices := []marketPrice.MarketPriceDto{
		{CurrentPrice: 165, PriceGets: 5, CurrentRate: 0.031},
	}
	mockMarketPriceRepo.On("FetchMarketPriceList", mock.Anything, []string{"AAPL"}).Return(mockMarketPrices, nil)

	mockDividend := &marketPrice.DividendEntity{DividendTotal: 2.0}
	mockMarketPriceRepo.On("FetchDividend", mock.Anything, "AAPL").Return(mockDividend, nil)

	// テスト対象メソッドの実行
	serviceInput := generated.UpdateUsStockInput{
		ID:        "1",
		GetPrice:  160,
		Quantity:  20,
		UsdJpy:    135.0,
	}
	updatedUsStock, err := service.UpdateUsStock(context.Background(), serviceInput)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUsStock)

	assert.Equal(t, "1", updatedUsStock.ID)
	assert.Equal(t, "AAPL", updatedUsStock.Code)
	assert.Equal(t, 160.0, updatedUsStock.GetPrice)
	assert.Equal(t, 20.0, updatedUsStock.Quantity)
	assert.Equal(t, "IT", updatedUsStock.Sector)
	assert.Equal(t, 135.0, updatedUsStock.UsdJpy)
	assert.Equal(t, 165.0, updatedUsStock.CurrentPrice)
	assert.Equal(t, 5.0, updatedUsStock.PriceGets)
	assert.Equal(t, 0.031, updatedUsStock.CurrentRate)
	assert.Equal(t, 2.0, updatedUsStock.Dividend)

	// モックの呼び出しを検証
	mockStockRepo.AssertExpectations(t)
	mockMarketPriceRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}


// TestDeleteUsStockService は DeleteUsStock メソッドのテストです。
func TestDeleteUsStockService(t *testing.T) {
	mockStockRepo := stock.NewMockUsStockRepository()
    mockMarketPriceRepo := marketPrice.NewMockMarketPriceRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewUsStockService(mockStockRepo, mockAuth, mockMarketPriceRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	// 成功時のテスト
	stockID := uint(1)
	mockStockRepo.On("DeleteUsStock", mock.Anything, stockID).Return(nil)

	// テスト対象メソッドの実行
	result, err := service.DeleteUsStock(context.Background(), "1")
	assert.NoError(t, err)
	assert.True(t, result)

	// モックの呼び出しを検証
	mockStockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
