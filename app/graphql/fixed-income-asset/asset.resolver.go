package fixedincomeasset

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    AssetService AssetService
}

func NewResolver(assetService AssetService) *Resolver {
    return &Resolver{AssetService: assetService}
}

func (r *Resolver) FixedIncomeAssets(ctx context.Context) ([]*generated.FixedIncomeAsset, error) {
    return r.AssetService.FixedIncomeAssets(ctx)
}

func (r *Resolver) CreateFixedIncomeAsset(ctx context.Context, input generated.CreateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
    newUsStock, err := r.AssetService.CreateFixedIncomeAsset(ctx, input)
    if err != nil {
        return nil, err
    }

    return newUsStock, nil
}