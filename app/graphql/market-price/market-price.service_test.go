package marketprice

import (
	"context"
	"errors"
	"my-us-stock-backend/app/graphql/generated"
	marketPrice "my-us-stock-backend/app/repository/market-price"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// FetchMarketPriceList のテスト
func TestFetchMarketPriceList(t *testing.T) {
    mockRepo := marketPrice.NewMockMarketPriceRepository()
    service := NewMarketPriceService(mockRepo)

    mockResponseBody := []marketPrice.MarketPriceDto{
        {Ticker: "AAPL", CurrentPrice: 189.84, PriceGets: 0.0685, CurrentRate: 0.13},
        {Ticker: "KO", CurrentPrice: 57.205, PriceGets: 0.0962, CurrentRate: 0.055},
    }
    mockResult := []*generated.MarketPrice{
        {Ticker: "AAPL", CurrentPrice: 189.84, PriceGets: 0.0685, CurrentRate: 0.13},
        {Ticker: "KO", CurrentPrice: 57.205, PriceGets: 0.0962, CurrentRate: 0.055},
    }
    tickers := []string{"AAPL", "KO"}
    mockRepo.On("FetchMarketPriceList", mock.Anything, tickers).Return(mockResponseBody, nil)

    result, err := service.FetchMarketPriceList(context.Background(), tickers)

    assert.NoError(t, err)
    assert.Equal(t, mockResult, result)

    mockRepo.AssertExpectations(t)
}

// エラー発生時のテスト
func TestFetchMarketPriceListError(t *testing.T) {
    mockRepo := marketPrice.NewMockMarketPriceRepository()
    service := NewMarketPriceService(mockRepo)

    tickers := []string{"INVALID"}
    mockRepo.On("FetchMarketPriceList", mock.Anything, tickers).Return([]marketPrice.MarketPriceDto(nil), errors.New("error fetching market prices"))

    _, err := service.FetchMarketPriceList(context.Background(), tickers)

    assert.Error(t, err)

    mockRepo.AssertExpectations(t)
}

