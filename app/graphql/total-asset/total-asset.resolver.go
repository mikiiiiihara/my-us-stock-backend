package totalasset

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    TotalAssetService TotalAssetService
}

func NewResolver(totalAssetService TotalAssetService) *Resolver {
    return &Resolver{TotalAssetService: totalAssetService}
}

func (r *Resolver) TotalAssets(ctx context.Context, day int) ([]*generated.TotalAsset, error) {
    return r.TotalAssetService.TotalAssets(ctx, day)
}