package main

import (
	"log"
	"my-us-stock-backend/graphql"
	"my-us-stock-backend/rest"

	"my-us-stock-backend/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    // .env ファイルから環境変数を読み込む
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    // PostgreSQLデータベースに接続
	db := database.Connect()
    
    // マイグレーションの実行
    database.Migrate(db)

    // Gin HTTPサーバーの初期化
    r := gin.Default() // gin.Engineのインスタンスを初期化

    // REST APIの設定
    rest.SetupREST(r, db)

    // GraphQLの設定
    graphql.SetupGraphQL(r, db)

    // サーバーを起動
    err = r.Run(":4000")
    if err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}