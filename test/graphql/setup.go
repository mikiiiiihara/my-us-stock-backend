package graphql

import (
	"context"
	serviceCurrency "my-us-stock-backend/app/graphql/currency"
	"my-us-stock-backend/app/graphql/generated"
	serviceMarketPrice "my-us-stock-backend/app/graphql/market-price"
	serviceUser "my-us-stock-backend/app/graphql/user"
)

type CustomQueryResolver struct {
    // 各リゾルバを含めます
    UserResolver         *serviceUser.Resolver
    CurrencyResolver     *serviceCurrency.Resolver
    MarketPriceResolver  *serviceMarketPrice.Resolver
}

// NewCustomQueryResolver - CustomQueryResolverのコンストラクタ関数
func NewCustomQueryResolver(userResolver *serviceUser.Resolver, currencyResolver *serviceCurrency.Resolver, marketPriceResolver *serviceMarketPrice.Resolver) *CustomQueryResolver {
    return &CustomQueryResolver{
        UserResolver:         userResolver,
        CurrencyResolver:     currencyResolver,
        MarketPriceResolver:  marketPriceResolver,
    }
}

func (r *CustomQueryResolver) Query() generated.QueryResolver {
    return r
}

func (r *CustomQueryResolver) Mutation() generated.MutationResolver {
    // ここでuserResolverを使用してMutationを実装する
    return r.UserResolver
}

func (r *CustomQueryResolver) User(ctx context.Context, id string) (*generated.User, error) {
    return r.UserResolver.User(ctx, id)
}

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