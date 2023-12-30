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

// UsStocks メソッドのテスト(0件の場合)
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