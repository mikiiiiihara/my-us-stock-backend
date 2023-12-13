// currency/service/currency.service.go
package currency

import (
	"context"
	"my-us-stock-backend/app/repository/currency"
)

type CurrencyService interface {
    FetchCurrentUsdJpy(ctx context.Context) (float64, error)
}

type DefaultCurrencyService struct {
    CurrencyRepo currency.CurrencyRepository // ポインタ型に変更
}

func NewCurrencyService(currencyRepo currency.CurrencyRepository) CurrencyService { // ポインタ型に変更
    return &DefaultCurrencyService{CurrencyRepo: currencyRepo}
}

func (s *DefaultCurrencyService) FetchCurrentUsdJpy(ctx context.Context) (float64, error) {
    return s.CurrencyRepo.FetchCurrentUsdJpy(ctx)
}
