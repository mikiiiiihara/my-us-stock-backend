package currency

import (
	"context"
)

type Resolver struct {
    CurrencyService CurrencyService
}

func NewResolver(currencyService CurrencyService) *Resolver {
    return &Resolver{CurrencyService: currencyService}
}

func (r *Resolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
    return r.CurrencyService.FetchCurrentUsdJpy(ctx)
}
