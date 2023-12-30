package marketprice

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockMarketPriceRepository の定義
type MockMarketPriceRepository struct {
    mock.Mock
}

// NewMockMarketPriceRepository は新しい MockMarketPriceRepository を作成し、初期設定を行います。
func NewMockMarketPriceRepository() *MockMarketPriceRepository {
	return &MockMarketPriceRepository{}
}

func (m *MockMarketPriceRepository) FetchMarketPriceList(ctx context.Context, tickers []string) ([]MarketPriceDto, error) {
    args := m.Called(ctx, tickers)
    return args.Get(0).([]MarketPriceDto), args.Error(1)
}

func (m *MockMarketPriceRepository) FetchDividend(ctx context.Context, ticker string) (*DividendEntity, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*DividendEntity), args.Error(1)
}