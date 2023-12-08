package graphql

import (
	"context"
	"my-us-stock-backend/graphql/currency"
	"my-us-stock-backend/graphql/generated"
	"my-us-stock-backend/graphql/user"
)

// CustomQueryResolver は QueryResolver と MutationResolver インターフェースを実装します
type CustomQueryResolver struct {
	UserResolver     *user.Resolver
	CurrencyResolver *currency.Resolver
}

// Queryメソッドの実装
func (r *CustomQueryResolver) Query() generated.QueryResolver {
	return r
}

// Userメソッドの実装
func (r *CustomQueryResolver) User(ctx context.Context, id string) (*generated.User, error) {
	return r.UserResolver.User(ctx, id)
}

// GetCurrentUsdJpyメソッドの実装
func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
	return r.CurrencyResolver.GetCurrentUsdJpy(ctx)
}

// Mutationメソッドの実装
func (r *CustomQueryResolver) Mutation() generated.MutationResolver {
	return r.UserResolver
}