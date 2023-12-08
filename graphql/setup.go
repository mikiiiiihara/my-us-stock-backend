package graphql

import (
	"my-us-stock-backend/graphql/currency"
	"my-us-stock-backend/graphql/generated"
	"my-us-stock-backend/graphql/user"

	repoCurrency "my-us-stock-backend/repository/currency"
	repoUser "my-us-stock-backend/repository/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler は GraphQL ハンドラをセットアップし、gin.HandlerFunc を返します
func Handler(userResolver *user.Resolver, currencyResolver *currency.Resolver) gin.HandlerFunc {
    resolver := &CustomQueryResolver{
        UserResolver:     userResolver,
        CurrencyResolver: currencyResolver,
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

    // GraphQLサービス、リゾルバの初期化
    currencyService := currency.NewCurrencyService(currencyRepo)
    currencyResolver := currency.NewResolver(currencyService)

    userService := user.NewUserService(userRepo)
    userResolver := user.NewResolver(userService)

    // GraphQL ハンドラ関数の設定
    r.POST("/graphql", Handler(userResolver, currencyResolver))
    r.GET("/graphql", PlaygroundHandler())
}
// Playgroundハンドラ関数
func PlaygroundHandler() gin.HandlerFunc {
    h := playground.Handler("GraphQL", "/graphql")

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}