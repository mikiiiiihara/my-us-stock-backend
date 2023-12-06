package currency

import (
	"context"
	"my-us-stock-backend/src/repository/currency"
)

type CurrencyService interface {
    FetchCurrentUsdJpy(ctx context.Context) (float64, error)
}

type DefaultCurrencyService struct {
    CurrencyRepo currency.CurrencyRepository
}

func NewCurrencyService(currencyRepo currency.CurrencyRepository) CurrencyService {
    return &DefaultCurrencyService{CurrencyRepo: currencyRepo}
}

func (s *DefaultCurrencyService) FetchCurrentUsdJpy(ctx context.Context) (float64, error) {
    return s.CurrencyRepo.FetchCurrentUsdJpy()
}
