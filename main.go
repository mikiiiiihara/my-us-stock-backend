package main

import (
	"context"
	"log"
	"my-us-stock-backend/graphql/currency"
	"my-us-stock-backend/graphql/generated"
	"my-us-stock-backend/graphql/user"
	repoCurrency "my-us-stock-backend/repository/currency"
	repoUser "my-us-stock-backend/repository/user"
	"my-us-stock-backend/repository/user/model"
	restUser "my-us-stock-backend/rest/user"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
    // モデルに基づいてテーブルを作成または更新
    db.AutoMigrate(&model.User{})
}

func main() {
        // .env ファイルから環境変数を読み込む
        err := godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file")
        }
    // PostgreSQLデータベースに接続
    dsn := "host=localhost user=myuser password=mypassword dbname=mydbname port=5432 sslmode=disable TimeZone=Asia/Tokyo"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

	// マイグレーションの実行
    Migrate(db)
    // リポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)
    currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    // GraphQLサービス、リゾルバの初期化
    currencyService := currency.NewCurrencyService(currencyRepo)
    currencyResolver := currency.NewResolver(currencyService)

    userService := user.NewUserService(userRepo)
    userResolver := user.NewResolver(userService)

    // RESTサービス、リゾルバの初期化
    userRestService := restUser.NewUserService(userRepo)
    userController := restUser.NewUserController(userRestService)
    // Gin HTTPサーバーの初期化
    r := gin.Default() // gin.Engineのインスタンスを初期化


    // RESTコントローラのルートを設定
    r.GET("/api/users/:id", userController.GetUser)
    r.POST("/api/users", userController.CreateUser)
    // GraphQL ハンドラ関数の設定
    r.POST("/graphql", graphqlHandler(userResolver, currencyResolver))
    r.GET("/graphql", playgroundHandler())

    // サーバーを起動
    err = r.Run(":4000")
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

// CustomQueryResolver は QueryResolver と MutationResolver インターフェースを実装します
type CustomQueryResolver struct {
    userResolver     *user.Resolver
    currencyResolver *currency.Resolver
}

// Queryメソッドの実装
func (r *CustomQueryResolver) Query() generated.QueryResolver {
    return r
}

// Userメソッドの実装
func (r *CustomQueryResolver) User(ctx context.Context, id string) (*generated.User, error) {
    return r.userResolver.User(ctx, id)
}

// GetCurrentUsdJpyメソッドの実装
func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
    return r.currencyResolver.GetCurrentUsdJpy(ctx)
}

// Mutationメソッドの実装
func (r *CustomQueryResolver) Mutation() generated.MutationResolver {
    return r.userResolver
}

// GraphQLハンドラ関数
func graphqlHandler(userResolver *user.Resolver, currencyResolver *currency.Resolver) gin.HandlerFunc {
    resolver := &CustomQueryResolver{
        userResolver:     userResolver,
        currencyResolver: currencyResolver,
    }
    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

    return func(c *gin.Context) {
        srv.ServeHTTP(c.Writer, c.Request)
    }
}



// Playgroundハンドラ関数
func playgroundHandler() gin.HandlerFunc {
    h := playground.Handler("GraphQL", "/graphql")

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
