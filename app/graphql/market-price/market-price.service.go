package marketprice

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
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
        return nil, utils.DefaultGraphQLError(err.Error())
    }

    marketPrices := make([]*generated.MarketPrice, len(dtos))
    for i, dto := range dtos {
        marketPrices[i] = &generated.MarketPrice{
            Ticker:       dto.Ticker,
            CurrentPrice: dto.CurrentPrice,
            PriceGets:    dto.PriceGets,
            CurrentRate:  dto.CurrentRate,
        }
    }    
    return marketPrices, nil
}
