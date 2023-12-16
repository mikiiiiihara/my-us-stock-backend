package marketprice

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    MarketPriceService MarketPriceService
}

func NewResolver(marketPriceService MarketPriceService) *Resolver {
    return &Resolver{MarketPriceService: marketPriceService}
}

func (r *Resolver) GetMarketPrices(ctx context.Context, tickers []string) ([]*generated.MarketPrice, error) {
    return r.MarketPriceService.FetchMarketPriceList(ctx, tickers)
}
