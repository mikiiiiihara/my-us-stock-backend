package marketprice

import (
	"context"
	"errors"
	"my-us-stock-backend/repository/market-price/dto"
	"my-us-stock-backend/repository/market-price/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMarketPriceRepository は CurrencyRepository のモックです。
type MockMarketPriceRepository struct {
    mock.Mock
}

func (m *MockMarketPriceRepository) FetchMarketPriceList(ctx context.Context, tickers []string) ([]dto.MarketPriceDto, error) {
	args := m.Called(ctx, tickers)
	return args.Get(0).([]dto.MarketPriceDto), args.Error(1)
}

func (m *MockMarketPriceRepository) FetchDividend(ctx context.Context, ticker string) (*entity.DividendEntity, error) {
	args := m.Called(ctx)
	return args.Get(0).(*entity.DividendEntity), args.Error(1)
}

func TestFetchMarketPriceList(t *testing.T) {
    mockRepo := new(MockMarketPriceRepository)
    service := NewMarketPriceService(mockRepo)

    mockResponseBody := []dto.MarketPriceDto{
        {Ticker: "AAPL", CurrentPrice: 189.84, PriceGets: 0.0685, CurrentRate: 0.13},
        {Ticker: "KO", CurrentPrice: 57.205, PriceGets: 0.0962, CurrentRate: 0.055},
    }
    tickers := []string{"AAPL", "KO"}
    mockRepo.On("FetchMarketPriceList", mock.Anything, tickers).Return(mockResponseBody, nil)

    result, err := service.FetchMarketPriceList(context.Background(), tickers)

    assert.NoError(t, err)
    assert.Equal(t, mockResponseBody, result)

    mockRepo.AssertExpectations(t)
}

func TestFetchMarketPriceListError(t *testing.T) {
    mockRepo := new(MockMarketPriceRepository)
    service := NewMarketPriceService(mockRepo)

    tickers := []string{"INVALID"}
    mockRepo.On("FetchMarketPriceList", mock.Anything, tickers).Return([]dto.MarketPriceDto(nil), errors.New("error fetching market prices"))

    _, err := service.FetchMarketPriceList(context.Background(), tickers)

    assert.Error(t, err)

    mockRepo.AssertExpectations(t)
}

