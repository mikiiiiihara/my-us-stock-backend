package marketprice

import (
	"context"
	"my-us-stock-backend/graphql/generated"
)

type Resolver struct {
    MarketPriceService MarketPriceService
}

func NewResolver(marketPriceService MarketPriceService) *Resolver {
    return &Resolver{MarketPriceService: marketPriceService}
}

func (r *Resolver) GetMarketPrices(ctx context.Context, tickers []string) ([]*generated.MarketPrice, error) {
    responses, err := r.MarketPriceService.FetchMarketPriceList(ctx, tickers)
    if err != nil {
        return nil, err
    }

    var marketPrices []*generated.MarketPrice
    for _, response := range responses {
        marketPrice := &generated.MarketPrice{
            Ticker:       response.Ticker,
            CurrentPrice: response.CurrentPrice,
            PriceGets:    response.PriceGets,
            CurrentRate:  response.CurrentRate,
        }
        marketPrices = append(marketPrices, marketPrice)
    }
    return marketPrices, nil
}
