package stock

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    UsStockService UsStockService
}

func NewResolver(usStockService UsStockService) *Resolver {
    return &Resolver{UsStockService: usStockService}
}

func (r *Resolver) UsStocks(ctx context.Context) ([]*generated.UsStock, error) {
    fmt.Println(r.UsStockService.UsStocks(ctx))
    return r.UsStockService.UsStocks(ctx)
}
