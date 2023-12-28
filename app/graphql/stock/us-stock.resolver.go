package stock

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    UsStockService UsStockService
}

func NewResolver(usStockService UsStockService) *Resolver {
    return &Resolver{UsStockService: usStockService}
}

func (r *Resolver) UsStocks(ctx context.Context) ([]*generated.UsStock, error) {
    return r.UsStockService.UsStocks(ctx)
}

func (r *Resolver) CreateUsStock(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error) {
    newUsStock, err := r.UsStockService.CreateUsStock(ctx, input)
    if err != nil {
        return nil, err
    }

    return newUsStock, nil
}