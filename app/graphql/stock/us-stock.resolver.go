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

func (r *Resolver) UpdateUsStock(ctx context.Context, input generated.UpdateUsStockInput) (*generated.UsStock, error) {
    newUsStock, err := r.UsStockService.UpdateUsStock(ctx, input)
    if err != nil {
        return nil, err
    }

    return newUsStock, nil
}

func (r *Resolver) DeleteUsStock(ctx context.Context, id string) (bool, error) {
    return r.UsStockService.DeleteUsStock(ctx, id)
}