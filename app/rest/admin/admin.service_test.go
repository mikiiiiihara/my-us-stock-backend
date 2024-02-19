package admin

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/market-price/fund"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test for FetchFundPrices method
func TestDefaultFundPriceService_FetchFundPrices(t *testing.T) {
	mockRepo := fund.NewMockFundPriceRepository()
	service := NewFundPriceService(mockRepo)

	expectedPrices := []model.FundPrice{
		{Name: "Fund A", Code: "FA123", Price: 100.50},
		{Name: "Fund B", Code: "FB456", Price: 200.75},
	}
	mockRepo.On("FetchFundPriceList", mock.Anything).Return(expectedPrices, nil)

	result, err := service.FetchFundPrices(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedPrices, result)
	mockRepo.AssertExpectations(t)
}

// Test for UpdateFundPrice method
func TestDefaultFundPriceService_UpdateFundPrice(t *testing.T) {
	mockRepo := fund.NewMockFundPriceRepository()
	service := NewFundPriceService(mockRepo)

	dto := fund.UpdateFundPriceDto{ID: 1, Price: 105.0}
	expectedPrice := &model.FundPrice{Name: "Fund A", Code: "FA123", Price: 105.0}

	mockRepo.On("UpdateFundPrice", mock.Anything, dto).Return(expectedPrice, nil)

	result, err := service.UpdateFundPrice(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, result)
	mockRepo.AssertExpectations(t)
}

// Test for CreateFundPrice method
func TestDefaultFundPriceService_CreateFundPrice(t *testing.T) {
	mockRepo := fund.NewMockFundPriceRepository()
	service := NewFundPriceService(mockRepo)

	dto := fund.CreateFundPriceDto{Name: "New Fund", Code: "NF789", Price: 250.0}
	expectedPrice := &model.FundPrice{Name: "New Fund", Code: "NF789", Price: 250.0}

	mockRepo.On("CreateFundPrice", mock.Anything, dto).Return(expectedPrice, nil)

	result, err := service.CreateFundPrice(context.Background(), dto)
	assert.NoError(t, err)
	assert.Equal(t, expectedPrice, result)
	mockRepo.AssertExpectations(t)
}
