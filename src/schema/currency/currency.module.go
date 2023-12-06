package currency

import (
	repoUser "my-us-stock-backend/src/repository/currency"
)

type CurrencyModule struct {
	CurrencyResolver *Resolver
}

func NewCurrencyModule() *CurrencyModule {
	currencyRepoModule := repoUser.NewCurrencyModule()
	currencyService := NewCurrencyService(*currencyRepoModule.Repository)
	currencyResolver := NewResolver(currencyService)

	return &CurrencyModule{
		CurrencyResolver: currencyResolver,
	}
}

func (um *CurrencyModule) Query() *Resolver {
	return um.CurrencyResolver
}

func (um *CurrencyModule) Mutation() *Resolver {
	return um.CurrencyResolver
}