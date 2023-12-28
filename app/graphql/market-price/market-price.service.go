package marketprice

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	repository "my-us-stock-backend/app/repository/market-price"
)

type MarketPriceService interface {
    FetchMarketPriceList(ctx context.Context, tickers []string) ([]*generated.MarketPrice, error)
}

type DefaultMarketPriceService struct {
    MarketPriceRepo repository.MarketPriceRepository
}

func NewMarketPriceService(marketPriceRepo repository.MarketPriceRepository) MarketPriceService {
    return &DefaultMarketPriceService{MarketPriceRepo: marketPriceRepo}
}

func (s *DefaultMarketPriceService) FetchMarketPriceList(ctx context.Context, tickers []string) ([]*generated.MarketPrice, error) {
    dtos, err := s.MarketPriceRepo.FetchMarketPriceList(ctx, tickers)
    if err != nil {
        return nil, err
    }

    marketPrices := make([]*generated.MarketPrice,0)
    for _, dto := range dtos {
        marketPrice := &generated.MarketPrice{
            Ticker:       dto.Ticker,
            CurrentPrice: dto.CurrentPrice,
            PriceGets:    dto.PriceGets,
            CurrentRate:  dto.CurrentRate,
        }
        marketPrices = append(marketPrices, marketPrice)
    }
    return marketPrices, nil
}
