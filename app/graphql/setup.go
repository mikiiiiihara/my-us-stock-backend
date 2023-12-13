package graphql

import (
	"my-us-stock-backend/app/graphql/currency"
	"my-us-stock-backend/app/graphql/generated"
	marketPrice "my-us-stock-backend/app/graphql/market-price"
	"my-us-stock-backend/app/graphql/user"

	repoCurrency "my-us-stock-backend/app/repository/currency"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoUser "my-us-stock-backend/app/repository/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler は GraphQL ハンドラをセットアップし、gin.HandlerFunc を返します
func Handler(userResolver *user.Resolver, currencyResolver *currency.Resolver,marketPriceResolver *marketPrice.Resolver) gin.HandlerFunc {
    resolver := &CustomQueryResolver{
        UserResolver:     userResolver,
        CurrencyResolver: currencyResolver,
        MarketPriceResolver: marketPriceResolver,
    }
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

    return func(c *gin.Context) {
        srv.ServeHTTP(c.Writer, c.Request)
    }
}

// SetupGraphQL は GraphQL ハンドラとリゾルバを設定します
func SetupGraphQL(r *gin.Engine, db *gorm.DB) {
    // リポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)
    currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    marketPriceRepo := repoMarketPrice.NewMarketPriceRepository(nil)

    // GraphQLサービス、リゾルバの初期化
    currencyService := currency.NewCurrencyService(currencyRepo)
    currencyResolver := currency.NewResolver(currencyService)

    marketPriceService := marketPrice.NewMarketPriceService(marketPriceRepo)
    marketPriceResolver := marketPrice.NewResolver(marketPriceService)

    userService := user.NewUserService(userRepo)
    userResolver := user.NewResolver(userService)

    // GraphQL ハンドラ関数の設定
    r.POST("/graphql", Handler(userResolver, currencyResolver, marketPriceResolver))
    r.GET("/graphql", PlaygroundHandler())
}
// Playgroundハンドラ関数
func PlaygroundHandler() gin.HandlerFunc {
    h := playground.Handler("GraphQL", "/graphql")

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}