package fixedincomeasset

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetService は MarketPriceService のモックです。
type MockAssetService struct {
    mock.Mock
}

func (m *MockAssetService) FixedIncomeAssets(ctx context.Context) ([]*generated.FixedIncomeAsset, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*generated.FixedIncomeAsset), args.Error(1)
}

func (m *MockAssetService) CreateFixedIncomeAsset(ctx context.Context, input generated.CreateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.FixedIncomeAsset), args.Error(1)
}

func (m *MockAssetService) UpdateFixedIncomeAsset(ctx context.Context, input generated.UpdateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.FixedIncomeAsset), args.Error(1)
}

func (m *MockAssetService) DeleteFixedIncomeAsset(ctx context.Context, id string) (bool, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(bool), args.Error(1)
}

// UsStocks メソッドのテスト
func TestFixedIncomeAssets(t *testing.T) {
    mockService := new(MockAssetService)
    resolver := NewResolver(mockService)

    fixedIncomeAssets := []*generated.FixedIncomeAsset{
        {ID: "1",Code: "i-Bond", GetPriceTotal: 10000.0, DividendRate: 1.5, PaymentMonth: []int{11}},
    }
    mockService.On("FixedIncomeAssets", mock.Anything).Return(fixedIncomeAssets, nil)

    result, err := resolver.FixedIncomeAssets(context.Background())
    
    assert.NoError(t, err)
    assert.Equal(t, fixedIncomeAssets, result)

    mockService.AssertExpectations(t)
}

// TestUpdateFixedIncomeAsset メソッドのテスト
func TestUpdateFixedIncomeAsset(t *testing.T) {
    mockService := new(MockAssetService)
    resolver := NewResolver(mockService)

    input := generated.UpdateFixedIncomeAssetInput{
        ID:            "1",
        GetPriceTotal: 15000.0,
    }
    mockResponse := &generated.FixedIncomeAsset{
        ID:            "1",
        Code:          "i-Bond",
        GetPriceTotal: 15000.0,
        DividendRate:  2.0,
        PaymentMonth:  []int{6, 12},
    }

    mockService.On("UpdateFixedIncomeAsset", mock.Anything, input).Return(mockResponse, nil)

    updatedAsset, err := resolver.UpdateFixedIncomeAsset(context.Background(), input)

    assert.NoError(t, err)
    assert.Equal(t, mockResponse, updatedAsset)

    mockService.AssertExpectations(t)
}


func TestCreateFixedIncomeAsset(t *testing.T) {
    mockService := new(MockAssetService)
    resolver := NewResolver(mockService)

	input := generated.CreateFixedIncomeAssetInput{
		Code: "Funds-からだにユーグレナファンド",
		GetPriceTotal: 100000.0, 
		DividendRate: 1.8,
		PaymentMonth: []int{3},

	}
	mockResponse := &generated.FixedIncomeAsset{
		Code: "Funds-からだにユーグレナファンド",
		GetPriceTotal: 100000.0, 
		DividendRate: 1.8,
		PaymentMonth: []int{3},
	}
    mockService.On("CreateFixedIncomeAsset", mock.Anything, input).Return(mockResponse, nil)

    result, err := resolver.CreateFixedIncomeAsset(context.Background(), input)
    
    assert.NoError(t, err)
    assert.Equal(t, mockResponse, result)

    mockService.AssertExpectations(t)
}

func TestDeleteFixedIncomeAsset(t *testing.T) {
    mockService := new(MockAssetService)
    resolver := NewResolver(mockService)

    mockService.On("DeleteFixedIncomeAsset", mock.Anything, "1").Return(true, nil)

    result, err := resolver.DeleteFixedIncomeAsset(context.Background(), "1")
    
    assert.NoError(t, err)
    assert.True(t, result)

    mockService.AssertExpectations(t)
}
