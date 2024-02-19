package admin

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/market-price/fund"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFundPriceService for testing
type MockFundPriceService struct {
	mock.Mock
}

func (m *MockFundPriceService) FetchFundPrices(ctx context.Context) ([]model.FundPrice, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.FundPrice), args.Error(1)
}

func (m *MockFundPriceService) UpdateFundPrice(ctx context.Context, dto fund.UpdateFundPriceDto) (*model.FundPrice, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.FundPrice), args.Error(1)
}

func (m *MockFundPriceService) CreateFundPrice(ctx context.Context, dto fund.CreateFundPriceDto) (*model.FundPrice, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.FundPrice), args.Error(1)
}

// Test for GetFundPrices method
func TestFundPriceController_GetFundPrices(t *testing.T) {
	mockService := new(MockFundPriceService)
	controller := NewFundPriceController(mockService)

	expectedPrices := []model.FundPrice{
		{Name: "Fund A", Code: "FA123", Price: 100.50},
		{Name: "Fund B", Code: "FB456", Price: 200.75},
	}
	mockService.On("FetchFundPrices", mock.Anything).Return(expectedPrices, nil)

	prices, err := controller.Service.FetchFundPrices(context.Background())

    assert.NoError(t, err)
    assert.Equal(t, expectedPrices, prices)
    mockService.AssertExpectations(t)
}

// Test for UpdateFundPrice method
func TestFundPriceController_UpdateFundPrice(t *testing.T) {
	mockService := new(MockFundPriceService)
	controller := NewFundPriceController(mockService)

	dto := fund.UpdateFundPriceDto{ID: 1, Price: 105.0}
	expectedPrice := &model.FundPrice{Name: "Fund A", Code: "FA123", Price: 105.0}

	mockService.On("UpdateFundPrice", mock.Anything, dto).Return(expectedPrice, nil)

	price, err := controller.Service.UpdateFundPrice(context.Background(), dto)

    assert.NoError(t, err)
    assert.Equal(t, expectedPrice, price)
    mockService.AssertExpectations(t)
}

// Test for CreateFundPrice method
func TestFundPriceController_CreateFundPrice(t *testing.T) {
	mockService := new(MockFundPriceService)
	controller := NewFundPriceController(mockService)

	dto := fund.CreateFundPriceDto{Name: "New Fund", Code: "NF789", Price: 250.0}
	expectedPrice := &model.FundPrice{Name: "New Fund", Code: "NF789", Price: 250.0}

	mockService.On("CreateFundPrice", mock.Anything, dto).Return(expectedPrice, nil)

	price, err := controller.Service.CreateFundPrice(context.Background(), dto)

    assert.NoError(t, err)
    assert.Equal(t, expectedPrice, price)
    mockService.AssertExpectations(t)
}
