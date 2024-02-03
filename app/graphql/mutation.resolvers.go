package graphql

import (
	"context"
	"my-us-stock-backend/app/graphql/crypto"
	FixedIncomeAsset "my-us-stock-backend/app/graphql/fixed-income-asset"
	"my-us-stock-backend/app/graphql/generated"
	JapanFund "my-us-stock-backend/app/graphql/japan-fund"
	"my-us-stock-backend/app/graphql/stock"
	TotalAsset "my-us-stock-backend/app/graphql/total-asset"
	"my-us-stock-backend/app/graphql/user"
)

// MutationResolver インターフェースを実装します
type CustomMutationResolver struct {
	UserResolver     *user.Resolver
	UsStockResolver *stock.Resolver
	CryptoResolver *crypto.Resolver
	FIxedIncomeAssetResolver *FixedIncomeAsset.Resolver
	JapanFundResolver *JapanFund.Resolver
	TotalAssetResolver *TotalAsset.Resolver
}

// Mutationメソッドの実装
func (r *CustomMutationResolver) Mutation() generated.MutationResolver {
	return r
}

func (r *CustomMutationResolver) CreateUser(ctx context.Context, input generated.CreateUserInput) (*generated.User, error) {
	return r.UserResolver.CreateUser(ctx, input)
}

func (r *CustomMutationResolver) CreateUsStock(ctx context.Context, input generated.CreateUsStockInput) (*generated.UsStock, error) {
	return r.UsStockResolver.CreateUsStock(ctx, input)
}

func (r *CustomMutationResolver) UpdateUsStock(ctx context.Context, input generated.UpdateUsStockInput) (*generated.UsStock, error) {
	return r.UsStockResolver.UpdateUsStock(ctx, input)
}

func (r *CustomMutationResolver) DeleteUsStock(ctx context.Context, id string) (bool, error) {
	return r.UsStockResolver.DeleteUsStock(ctx, id)
}

func (r *CustomMutationResolver) CreateCrypto(ctx context.Context, input generated.CreateCryptoInput) (*generated.Crypto, error) {
	return r.CryptoResolver.CreateCrypto(ctx, input)
}

func (r *CustomMutationResolver) UpdateCrypto(ctx context.Context, input generated.UpdateCryptoInput) (*generated.Crypto, error) {
	return r.CryptoResolver.UpdateCrypto(ctx, input)
}

func (r *CustomMutationResolver) DeleteCrypto(ctx context.Context, id string) (bool, error) {
	return r.CryptoResolver.DeleteCrypto(ctx, id)
}

func (r *CustomMutationResolver) CreateFixedIncomeAsset(ctx context.Context, input generated.CreateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
	return r.FIxedIncomeAssetResolver.CreateFixedIncomeAsset(ctx, input)
}

func (r *CustomMutationResolver) UpdateFixedIncomeAsset(ctx context.Context, input generated.UpdateFixedIncomeAssetInput) (*generated.FixedIncomeAsset, error) {
	return r.FIxedIncomeAssetResolver.UpdateFixedIncomeAsset(ctx, input)
}

func (r *CustomMutationResolver) DeleteFixedIncomeAsset(ctx context.Context, id string) (bool, error) {
	return r.FIxedIncomeAssetResolver.DeleteFixedIncomeAsset(ctx, id)
}

func (r *CustomMutationResolver) CreateJapanFund(ctx context.Context, input generated.CreateJapanFundInput) (*generated.JapanFund, error) {
	return r.JapanFundResolver.CreateJapanFund(ctx, input)
}

func (r *CustomMutationResolver) UpdateJapanFund(ctx context.Context, input generated.UpdateJapanFundInput) (*generated.JapanFund, error) {
	return r.JapanFundResolver.UpdateJapanFund(ctx, input)
}

func (r *CustomMutationResolver) DeleteJapanFund(ctx context.Context, id string) (bool, error) {
	return r.JapanFundResolver.DeleteJapanFund(ctx, id)
}

func (r *CustomMutationResolver) UpdateTotalAsset(ctx context.Context, input generated.UpdateTotalAssetInput) (*generated.TotalAsset, error) {
	return r.TotalAssetResolver.UpdateTotalAsset(ctx, input)
}