package graphql

import (
	"context"
	"my-us-stock-backend/app/graphql/crypto"
	"my-us-stock-backend/app/graphql/currency"
	FixedIncomeAsset "my-us-stock-backend/app/graphql/fixed-income-asset"
	"my-us-stock-backend/app/graphql/generated"
	JapanFund "my-us-stock-backend/app/graphql/japan-fund"
	marketPrice "my-us-stock-backend/app/graphql/market-price"
	"my-us-stock-backend/app/graphql/stock"
	"my-us-stock-backend/app/graphql/user"
)

// QueryResolverインターフェースを実装します
type CustomQueryResolver struct {
	UserResolver     *user.Resolver
	CurrencyResolver *currency.Resolver
	MarketPriceResolver *marketPrice.Resolver
	UsStockResolver *stock.Resolver
	CryptoResolver *crypto.Resolver
	FIxedIncomeAssetResolver *FixedIncomeAsset.Resolver
	JapanFundResolver *JapanFund.Resolver
}

// Queryメソッドの実装
func (r *CustomQueryResolver) Query() generated.QueryResolver {
	return r
}

func (r *CustomQueryResolver) User(ctx context.Context) (*generated.User, error) {
	return r.UserResolver.User(ctx)
}

func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
	return r.CurrencyResolver.GetCurrentUsdJpy(ctx)
}

func (r *CustomQueryResolver) GetMarketPrices(ctx context.Context, tickers []*string) ([]*generated.MarketPrice, error) {
    // 文字列スライスに変換
    tickerStrs := make([]string, len(tickers))
    for i, t := range tickers {
        tickerStrs[i] = *t
    }

    // サービスを呼び出して結果を取得
    return r.MarketPriceResolver.GetMarketPrices(ctx, tickerStrs)
}

func (r *CustomQueryResolver) UsStocks(ctx context.Context) ([]*generated.UsStock, error) {
	return r.UsStockResolver.UsStocks(ctx)
}

func (r *CustomQueryResolver) Cryptos(ctx context.Context) ([]*generated.Crypto, error) {
	return r.CryptoResolver.Cryptos(ctx)
}

func (r *CustomQueryResolver) FixedIncomeAssets(ctx context.Context) ([]*generated.FixedIncomeAsset, error) {
	return r.FIxedIncomeAssetResolver.FixedIncomeAssets(ctx)
}

func (r *CustomQueryResolver) JapanFunds(ctx context.Context) ([]*generated.JapanFund, error) {
	return r.JapanFundResolver.JapanFunds (ctx)
}
