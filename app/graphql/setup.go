package graphql

import (
	"context"
	authService "my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/common/auth/logic"
	"my-us-stock-backend/app/common/auth/validation"
	"my-us-stock-backend/app/graphql/currency"
	"my-us-stock-backend/app/graphql/generated"
	marketPrice "my-us-stock-backend/app/graphql/market-price"
	"my-us-stock-backend/app/graphql/stock"
	"my-us-stock-backend/app/graphql/user"
	"my-us-stock-backend/app/graphql/utils"

	repoStock "my-us-stock-backend/app/repository/assets/stock"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoCurrency "my-us-stock-backend/app/repository/market-price/currency"
	repoUser "my-us-stock-backend/app/repository/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// この関数は、GinのContextからGraphQLのContextにデータを転送します。
func ginContextToGraphQLMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Cookieの取得、見つからない場合は空文字とする
        cookie, _ := c.Cookie("access_token")

        // GraphQLのContextにCookieを追加（空文字も含む）
        ctx := context.WithValue(c.Request.Context(), utils.CookieKey, cookie)
        c.Request = c.Request.WithContext(ctx)

        c.Next()
    }
}


// Handler は GraphQL ハンドラをセットアップし、gin.HandlerFunc を返します
func Handler(userResolver *user.Resolver, currencyResolver *currency.Resolver,marketPriceResolver *marketPrice.Resolver, usStockResolver *stock.Resolver) gin.HandlerFunc {
    resolver := &CustomQueryResolver{
        UserResolver:     userResolver,
        CurrencyResolver: currencyResolver,
        MarketPriceResolver: marketPriceResolver,
        UsStockResolver: usStockResolver,
    }
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

    return func(c *gin.Context) {
        // ここでGinのContextをGraphQLのContextに変換
        ginContextToGraphQLMiddleware()(c)

        srv.ServeHTTP(c.Writer, c.Request)
    }
}

// SetupGraphQL は GraphQL ハンドラとリゾルバを設定します
func SetupGraphQL(r *gin.Engine, db *gorm.DB) {
    // リポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)
    currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    marketPriceRepo := repoMarketPrice.NewMarketPriceRepository(nil)
    usStockRepo := repoStock.NewUsStockRepository(db)

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
    
    usStockService := stock.NewUsStockService(usStockRepo, authService)
    usStockResolver := stock.NewResolver(usStockService)
    // GraphQLエンドポイントへのルート設定
    r.POST("/graphql", ginContextToGraphQLMiddleware(), Handler(userResolver, currencyResolver, marketPriceResolver,usStockResolver))
    r.GET("/graphql", PlaygroundHandler())
}
// Playgroundハンドラ関数
func PlaygroundHandler() gin.HandlerFunc {
    h := playground.Handler("GraphQL", "/graphql")

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}