package graphql

import (
	"context"
	"my-us-stock-backend/app/graphql/utils"

	"github.com/gin-gonic/gin"
)

// この関数は、GinのContextからGraphQLのContextにデータを転送します。
func GinContextToGraphQLMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Cookieの取得、見つからない場合は空文字とする
        cookie, _ := c.Cookie("access_token")

        // GraphQLのContextにCookieを追加（空文字も含む）
        ctx := context.WithValue(c.Request.Context(), utils.CookieKey, cookie)
        c.Request = c.Request.WithContext(ctx)

        c.Next()
    }
}