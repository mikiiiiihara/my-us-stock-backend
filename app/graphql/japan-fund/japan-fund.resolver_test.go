package japanfund

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJapanFundService は MarketPriceService のモックです。
type MockJapanFundService struct {
    mock.Mock
}

func (m *MockJapanFundService) JapanFunds(ctx context.Context) ([]*generated.JapanFund, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*generated.JapanFund), args.Error(1)
}

func (m *MockJapanFundService) CreateJapanFund(ctx context.Context, input generated.CreateJapanFundInput) (*generated.JapanFund, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.JapanFund), args.Error(1)
}

// UsStocks メソッドのテスト
func TestJapanFunds(t *testing.T) {
    mockService := new(MockJapanFundService)
    resolver := NewResolver(mockService)

    japanFunds := []*generated.JapanFund{
        {ID: "1",Code: "SP500", Name:"ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", GetPrice: 15523.81, GetPriceTotal: 761157,CurrentPrice: 24281},
    }
    mockService.On("JapanFunds", mock.Anything).Return(japanFunds, nil)

    result, err := resolver.JapanFunds(context.Background())
    
    assert.NoError(t, err)
    assert.Equal(t, japanFunds, result)

    mockService.AssertExpectations(t)
}

// UsStocks メソッドのテスト(0件の場合)
func TestCreateJapanFund(t *testing.T) {
    mockService := new(MockJapanFundService)
    resolver := NewResolver(mockService)

	input := generated.CreateJapanFundInput{
		Code: "全世界株",
		Name: "ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）",
		GetPrice: 18609,
		GetPriceTotal: 400004,

	}
	mockResponse := &generated.JapanFund{
		ID: "1",
		Code: "全世界株",
		Name: "ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）",
		GetPrice: 18609,
		GetPriceTotal: 400004,
		CurrentPrice: 21084,
	}
    mockService.On("CreateJapanFund", mock.Anything, input).Return(mockResponse, nil)

    result, err := resolver.CreateJapanFund(context.Background(), input)
    
    assert.NoError(t, err)
    assert.Equal(t, mockResponse, result)

    mockService.AssertExpectations(t)
}