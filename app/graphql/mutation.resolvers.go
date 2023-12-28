package graphql

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/stock"
	"my-us-stock-backend/app/graphql/user"
)

// MutationResolver インターフェースを実装します
type CustomMutationResolver struct {
	UserResolver     *user.Resolver
	UsStockResolver *stock.Resolver
}

// Mutationメソッドの実装
func (r *CustomMutationResolver) Mutation() generated.MutationResolver {
	return r
}

// Userメソッドの実装
func (r *CustomMutationResolver) CreateUser(ctx context.Context, input generated.CreateUserInput) (*generated.User, error) {
	return r.UserResolver.CreateUser(ctx, input)
}

// UsStockメソッドの実装
func (r *CustomMutationResolver) CreateUsStock(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error) {
	return r.UsStockResolver.CreateUsStock(ctx, input)
}

