package stock

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsStockService は MarketPriceService のモックです。
type MockUsStockService struct {
    mock.Mock
}

func (m *MockUsStockService) UsStocks(ctx context.Context) ([]*generated.UsStock, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*generated.UsStock), args.Error(1)
}

func (m *MockUsStockService) CreateUsStock(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.UsStock), args.Error(1)
}

func (m *MockUsStockService) UpdateUsStock(ctx context.Context, input generated.UpdateUsStockInput) (*generated.UsStock, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.UsStock), args.Error(1)
}

func (m *MockUsStockService) DeleteUsStock(ctx context.Context, id string) (bool, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(bool), args.Error(1)
}

// UsStocks メソッドのテスト
func TestUsStocks(t *testing.T) {
    mockService := new(MockUsStockService)
    resolver := NewResolver(mockService)

    usStocks := []*generated.UsStock{
        {ID: "1",Code: "AAPL", GetPrice: 180.0, Dividend: 1.22, Quantity: 2, Sector: "IT", UsdJpy: 130.2,CurrentPrice: 189.84, PriceGets: 0.0685, CurrentRate: 0.13},
        {ID: "2",Code: "KO",  GetPrice: 50.0, Dividend: 1.22, Quantity: 2, Sector: "Consumer Staples", UsdJpy: 130.2,CurrentPrice: 57.205, PriceGets: 0.0962, CurrentRate: 0.055},
    }
    mockService.On("UsStocks", mock.Anything).Return(usStocks, nil)

    result, err := resolver.UsStocks(context.Background())
    
    assert.NoError(t, err)
    assert.Equal(t, usStocks, result)

    mockService.AssertExpectations(t)
}

// UsStocks メソッドのテスト(0件の場合)
func TestCreateUsStock(t *testing.T) {
    mockService := new(MockUsStockService)
    resolver := NewResolver(mockService)

	input := generated.CreateUsStockInput{
		Code: "AAPL",
		GetPrice: 180.0, 
		Quantity: 2, 
		Sector: "IT", 
		UsdJpy: 130.2,

	}
	mockResponse := &generated.UsStock{
		ID: "1",
		Code: "AAPL",
		GetPrice: 180.0, 
		Quantity: 2, 
		Sector: "IT", 
		UsdJpy: 130.2,
		Dividend: 1.22,
		CurrentPrice: 189.84,
		PriceGets: 0.0685, 
		CurrentRate: 0.13,
	}
    mockService.On("CreateUsStock", mock.Anything, input).Return(mockResponse, nil)

    result, err := resolver.CreateUsStock(context.Background(), input)
    
    assert.NoError(t, err)
    assert.Equal(t, mockResponse, result)

    mockService.AssertExpectations(t)
}

// TestUpdateUsStock メソッドのテスト
func TestUpdateUsStock(t *testing.T) {
    mockService := new(MockUsStockService)
    resolver := NewResolver(mockService)

    input := generated.UpdateUsStockInput{
        ID: "1",
        GetPrice: 190.0, 
        Quantity: 3, 
        UsdJpy: 131.0,
    }

    mockResponse := &generated.UsStock{
        ID: "1",
        Code: "AAPL",
        GetPrice: 190.0, 
        Quantity: 3, 
        Sector: "IT", 
        UsdJpy: 131.0,
        Dividend: 1.24,
        CurrentPrice: 192.84,
        PriceGets: 0.0725, 
        CurrentRate: 0.14,
    }

    mockService.On("UpdateUsStock", mock.Anything, input).Return(mockResponse, nil)

    result, err := resolver.UpdateUsStock(context.Background(), input)

    assert.NoError(t, err)
    assert.Equal(t, mockResponse, result)

    mockService.AssertExpectations(t)
}


func TestDeleteUsStock(t *testing.T) {
    mockService := new(MockUsStockService)
    resolver := NewResolver(mockService)

    mockService.On("DeleteUsStock", mock.Anything, "1").Return(true, nil)

    result, err := resolver.DeleteUsStock(context.Background(), "1")
    
    assert.NoError(t, err)
    assert.True(t, result)

    mockService.AssertExpectations(t)
}
