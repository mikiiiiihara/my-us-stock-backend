package graphql

import (
	authService "my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/common/auth/logic"
	"my-us-stock-backend/app/common/auth/validation"
	"my-us-stock-backend/app/graphql/crypto"
	"my-us-stock-backend/app/graphql/currency"
	fixedIncomeAsset "my-us-stock-backend/app/graphql/fixed-income-asset"
	"my-us-stock-backend/app/graphql/generated"
	japanFund "my-us-stock-backend/app/graphql/japan-fund"
	marketPrice "my-us-stock-backend/app/graphql/market-price"
	"my-us-stock-backend/app/graphql/stock"
	totalAsset "my-us-stock-backend/app/graphql/total-asset"
	"my-us-stock-backend/app/graphql/user"

	repoCrypto "my-us-stock-backend/app/repository/assets/crypto"
	repoFixedIncome "my-us-stock-backend/app/repository/assets/fixed-income"
	repoJapanFund "my-us-stock-backend/app/repository/assets/fund"
	repoStock "my-us-stock-backend/app/repository/assets/stock"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoMarketCrypto "my-us-stock-backend/app/repository/market-price/crypto"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoTotalAsset "my-us-stock-backend/app/repository/total-assets"
	repoUser "my-us-stock-backend/app/repository/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CombinedResolverは、クエリとミューテーションの両方を処理するリゾルバです。
type CombinedResolver struct {
    *CustomQueryResolver
    *CustomMutationResolver
}


// Handler は GraphQL ハンドラをセットアップし、gin.HandlerFunc を返します
func Handler(userResolver *user.Resolver, currencyResolver *currency.Resolver,marketPriceResolver *marketPrice.Resolver, usStockResolver *stock.Resolver, cryptoResolver *crypto.Resolver, fixedIncomeAssetResolver *fixedIncomeAsset.Resolver, japanFundResolver *japanFund.Resolver, totalAssetResolver *totalAsset.Resolver) gin.HandlerFunc {
    queryResolver := &CustomQueryResolver{
        UserResolver:     userResolver,
        CurrencyResolver: currencyResolver,
        MarketPriceResolver: marketPriceResolver,
        UsStockResolver: usStockResolver,
        CryptoResolver: cryptoResolver,
        FIxedIncomeAssetResolver: fixedIncomeAssetResolver,
        JapanFundResolver: japanFundResolver,
        TotalAssetResolver: totalAssetResolver,
    }
    mutationResolver := &CustomMutationResolver{
        UserResolver: userResolver,
        UsStockResolver: usStockResolver,
        CryptoResolver: cryptoResolver,
        FIxedIncomeAssetResolver: fixedIncomeAssetResolver,
        JapanFundResolver: japanFundResolver,
    }
    combinedResolver := &CombinedResolver{
        CustomQueryResolver: queryResolver,
        CustomMutationResolver: mutationResolver,
    }
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: combinedResolver}))

    return func(c *gin.Context) {
        // ここでGinのContextをGraphQLのContextに変換
        GinContextToGraphQLMiddleware()(c)

        srv.ServeHTTP(c.Writer, c.Request)
    }
}

// SetupGraphQL は GraphQL ハンドラとリゾルバを設定します
func SetupGraphQL(r *gin.Engine, db *gorm.DB) {
    // リポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)
    currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    marketPriceRepo := repoMarketPrice.NewMarketPriceRepository(nil)
    marketCryptoRepo := repoMarketCrypto.NewCryptoRepository(nil)
    usStockRepo := repoStock.NewUsStockRepository(db)
    cryptoRepo := repoCrypto.NewCryptoRepository(db)
    japanFundRepo := repoJapanFund.NewJapanFundRepository(db)
    fixedIncomeAssetRepo := repoFixedIncome.NewFixedIncomeRepository(db)
    totalAssetRepo := repoTotalAsset.NewTotalAssetRepository(db)

    // 認証機能
    userLogic := logic.NewUserLogic()
    responseLogic := logic.NewResponseLogic()
    jwtLogic := logic.NewJWTLogic()
    authValidation := validation.NewAuthValidation()

    // 認証サービスの初期化
    authService := authService.NewAuthService(userRepo, userLogic, responseLogic, jwtLogic, authValidation)

    // GraphQLサービス、リゾルバの初期化
    currencyService := currency.NewCurrencyService(currencyRepo)
    currencyResolver := currency.NewResolver(currencyService)

    marketPriceService := marketPrice.NewMarketPriceService(marketPriceRepo)
    marketPriceResolver := marketPrice.NewResolver(marketPriceService)

    userService := user.NewUserService(userRepo,authService)
    userResolver := user.NewResolver(userService)
    
    usStockService := stock.NewUsStockService(usStockRepo, authService, marketPriceRepo)
    usStockResolver := stock.NewResolver(usStockService)

    cryptoService := crypto.NewCryptoService(cryptoRepo, authService, marketCryptoRepo)
    cryptoResolver := crypto.NewResolver(cryptoService)

    fixedIncomeAssetService := fixedIncomeAsset.NewAssetService(fixedIncomeAssetRepo, authService)
    fixedIncomeAssetResolver := fixedIncomeAsset.NewResolver(fixedIncomeAssetService)

    japanFundService := japanFund.NewJapanFundService(japanFundRepo, authService)
    japanFundResolver := japanFund.NewResolver(japanFundService)

    totalAssetService := totalAsset.NewTotalAssetService(totalAssetRepo, authService)
    totalAssetResolver := totalAsset.NewResolver(totalAssetService)

    // GraphQLエンドポイントへのルート設定
    r.POST("/graphql", GinContextToGraphQLMiddleware(), Handler(userResolver, currencyResolver, marketPriceResolver,usStockResolver, cryptoResolver,fixedIncomeAssetResolver, japanFundResolver,totalAssetResolver))
    r.GET("/graphql", PlaygroundHandler())
}
// Playgroundハンドラ関数
func PlaygroundHandler() gin.HandlerFunc {
    h := playground.Handler("GraphQL", "/graphql")

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}