package main

import (
	"log"
	"my-us-stock-backend/src/repository/user/model"
	"my-us-stock-backend/src/schema"
	"my-us-stock-backend/src/schema/generated"

	"my-us-stock-backend/src/controller"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
    // モデルに基づいてテーブルを作成または更新
    db.AutoMigrate(&model.User{})
}

func main() {
    // PostgreSQLデータベースに接続
    dsn := "host=localhost user=myuser password=mypassword dbname=mydbname port=5432 sslmode=disable TimeZone=Asia/Tokyo"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

	// マイグレーションの実行
    Migrate(db)

    // コントローラレジストリの作成
    controllerModule := controller.NewControllerModule(db)

    // Gin HTTPサーバーの初期化
    r := gin.Default()

    // コントローラレジストリを使用してREST APIルートを登録
    controllerModule.RegisterRoutes(r)

    // GraphQLのエンドポイントのセットアップ
    r.POST("/graphql", graphqlHandler(db))
    r.GET("/graphql", playgroundHandler())

    // サーバーを起動
    err = r.Run(":4000")
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

// GraphQLハンドラ関数
func graphqlHandler(db *gorm.DB) gin.HandlerFunc {
	resolver := schema.NewSchemaModule(db)
    h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}

// Playgroundハンドラ関数
func playgroundHandler() gin.HandlerFunc {
    h := playground.Handler("GraphQL", "/graphql")

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
