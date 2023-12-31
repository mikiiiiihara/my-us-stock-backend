package marketprice

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMarketPriceService は MarketPriceService のモックです。
type MockMarketPriceService struct {
    mock.Mock
}

func (m *MockMarketPriceService) FetchMarketPriceList(ctx context.Context, tickers []string) ([]*generated.MarketPrice, error) {
    args := m.Called(ctx, tickers)
    return args.Get(0).([]*generated.MarketPrice), args.Error(1)
}

// GetMarketPrices メソッドのテスト
func TestGetMarketPrices(t *testing.T) {
    mockService := new(MockMarketPriceService)
    resolver := NewResolver(mockService)

    mockMarketPrices := []*generated.MarketPrice{
        {Ticker: "AAPL", CurrentPrice: 189.84, PriceGets: 0.0685, CurrentRate: 0.13},
        {Ticker: "KO", CurrentPrice: 57.205, PriceGets: 0.0962, CurrentRate: 0.055},
    }
    tickers := []string{"AAPL", "KO"}
    mockService.On("FetchMarketPriceList", mock.Anything, tickers).Return(mockMarketPrices, nil)

    result, err := resolver.MarketPrices(context.Background(), tickers)
    
    assert.NoError(t, err)
    assert.Equal(t, mockMarketPrices, result)

    mockService.AssertExpectations(t)
}