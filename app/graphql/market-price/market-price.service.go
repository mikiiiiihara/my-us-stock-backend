package marketprice

import (
	"context"
	repository "my-us-stock-backend/app/repository/market-price"
	"my-us-stock-backend/app/repository/market-price/dto"
)

type MarketPriceService interface {
    FetchMarketPriceList(ctx context.Context, tickers []string) ([]dto.MarketPriceDto, error)
}

type DefaultMarketPriceService struct {
    MarketPriceRepo repository.MarketPriceRepository // ポインタ型に変更
}

func NewMarketPriceService(marketPriceRepo repository.MarketPriceRepository) MarketPriceService { // ポインタ型に変更
    return &DefaultMarketPriceService{MarketPriceRepo: marketPriceRepo}
}

func (s *DefaultMarketPriceService) FetchMarketPriceList(ctx context.Context, tickers []string) ([]dto.MarketPriceDto, error) {
    return s.MarketPriceRepo.FetchMarketPriceList(ctx,tickers);
}