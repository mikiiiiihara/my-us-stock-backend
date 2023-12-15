package graphql

import (
	"context"
	serviceCurrency "my-us-stock-backend/app/graphql/currency"
	"my-us-stock-backend/app/graphql/generated"
	serviceMarketPrice "my-us-stock-backend/app/graphql/market-price"
	serviceUser "my-us-stock-backend/app/graphql/user"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoUser "my-us-stock-backend/app/repository/user"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"gorm.io/gorm"
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

// SetupOptions - GraphQLサーバーのセットアップオプション
type SetupOptions struct {
    MockHTTPClient    *http.Client
    CurrencyRepo      repoCurrency.CurrencyRepository
    UserRepo          repoUser.UserRepository
    MarketPriceRepo   repoMarketPrice.MarketPriceRepository
}

// SetupGraphQLServer - GraphQLサーバーのセットアップ
func SetupGraphQLServer(db *gorm.DB, opts *SetupOptions) *handler.Server {
    var currencyRepo repoCurrency.CurrencyRepository
    var userRepo repoUser.UserRepository
    var marketPriceRepo repoMarketPrice.MarketPriceRepository

    // optsがnilでない場合にのみ、各リポジトリを設定
    if opts != nil {
        currencyRepo = opts.CurrencyRepo
        userRepo = opts.UserRepo
        marketPriceRepo = opts.MarketPriceRepo
    }

    // 各リポジトリがまだnilの場合、デフォルトのリポジトリを使用
    if currencyRepo == nil {
        currencyRepo = repoCurrency.NewCurrencyRepository(nil)
    }
    if userRepo == nil {
        userRepo = repoUser.NewUserRepository(db)
    }
    if marketPriceRepo == nil {
        // 注意: ここでは opts が nil の可能性があるため、opts.MockHTTPClient の前に nil チェックが必要
        var httpClient *http.Client
        if opts != nil {
            httpClient = opts.MockHTTPClient
        }
        marketPriceRepo = repoMarketPrice.NewMarketPriceRepository(httpClient)
    }

    // サービスとリゾルバの初期化
    currencyService := serviceCurrency.NewCurrencyService(currencyRepo)
    currencyResolver := serviceCurrency.NewResolver(currencyService)

    userService := serviceUser.NewUserService(userRepo)
    userResolver := serviceUser.NewResolver(userService)

    marketPriceService := serviceMarketPrice.NewMarketPriceService(marketPriceRepo)
    marketPriceResolver := serviceMarketPrice.NewResolver(marketPriceService)

    // CustomQueryResolverの初期化
    resolver := NewCustomQueryResolver(userResolver, currencyResolver, marketPriceResolver)

    return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
}
