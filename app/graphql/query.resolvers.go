package graphql

import (
	"context"
	"my-us-stock-backend/app/graphql/currency"
	"my-us-stock-backend/app/graphql/generated"
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
}

// Queryメソッドの実装
func (r *CustomQueryResolver) Query() generated.QueryResolver {
	return r
}

// Userメソッドの実装
func (r *CustomQueryResolver) User(ctx context.Context) (*generated.User, error) {
	return r.UserResolver.User(ctx)
}

// GetCurrentUsdJpyメソッドの実装
func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
	return r.CurrencyResolver.GetCurrentUsdJpy(ctx)
}

// GetMarketPricesメソッドの実装
func (r *CustomQueryResolver) GetMarketPrices(ctx context.Context, tickers []*string) ([]*generated.MarketPrice, error) {
    // 文字列スライスに変換
    tickerStrs := make([]string, len(tickers))
    for i, t := range tickers {
        tickerStrs[i] = *t
    }

    // サービスを呼び出して結果を取得
    return r.MarketPriceResolver.GetMarketPrices(ctx, tickerStrs)
}

// UsStocksメソッドの実装
func (r *CustomQueryResolver) UsStocks(ctx context.Context) ([]*generated.UsStock, error) {
	return r.UsStockResolver.UsStocks(ctx)
}