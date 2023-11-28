package main

import (
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/playground"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
    "my-us-stock-backend/graph/generated"
	"my-us-stock-backend/graph"
	"my-us-stock-backend/graph/user/model"
)

func Migrate(db *gorm.DB) {
    // モデルに基づいてテーブルを作成または更新
    db.AutoMigrate(&model.User{})
}

func main() {
    // PostgreSQLデータベースに接続
    dsn := "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Tokyo"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

	// マイグレーションの実行
    Migrate(db)

    // Gin HTTPサーバーの初期化
    r := gin.Default()

    // GraphQLのエンドポイントのセットアップ
    r.POST("/graphql", graphqlHandler(db))
    r.GET("/", playgroundHandler())

    // サーバーを起動
    err = r.Run() // デフォルトでは ":8080" でサーバーを起動
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

// GraphQLハンドラ関数
func graphqlHandler(db *gorm.DB) gin.HandlerFunc {
	resolver := graph.NewResolver(db)
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
