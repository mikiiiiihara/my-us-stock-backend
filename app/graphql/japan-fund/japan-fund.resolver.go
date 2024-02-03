package japanfund

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    JapanFundService JapanFundService
}

func NewResolver(usStockService JapanFundService) *Resolver {
    return &Resolver{JapanFundService: usStockService}
}

func (r *Resolver) JapanFunds(ctx context.Context) ([]*generated.JapanFund, error) {
    return r.JapanFundService.JapanFunds(ctx)
}

func (r *Resolver) CreateJapanFund(ctx context.Context, input generated.CreateJapanFundInput) (*generated.JapanFund, error) {
    newFund, err := r.JapanFundService.CreateJapanFund(ctx, input)
    if err != nil {
        return nil, err
    }

    return newFund, nil
}

func (r *Resolver) UpdateJapanFund(ctx context.Context, input generated.UpdateJapanFundInput) (*generated.JapanFund, error) {
    newFund, err := r.JapanFundService.UpdateJapanFund(ctx, input)
    if err != nil {
        return nil, err
    }

    return newFund, nil
}

func (r *Resolver) DeleteJapanFund(ctx context.Context, id string) (bool, error) {
    return r.JapanFundService.DeleteJapanFund(ctx, id)
}