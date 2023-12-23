package auth

import (
	"net/http"
	"strings"

	"my-us-stock-backend/app/common/auth"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService auth.AuthService
}

func NewAuthController(authService auth.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// SignIn ログイン
func (ac *AuthController) SignIn(c *gin.Context) {
    ctx := c.Request.Context() // context.Context を取得
    user, err := ac.AuthService.SignIn(ctx, c)
    if err != nil {
        // エラーハンドリングをここに追加する必要があります。
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ac.AuthService.SendAuthResponse(ctx, c, user, http.StatusOK)
}


// SignUp 会員登録処理
func (ac *AuthController) SignUp(c *gin.Context) {
	ctx := c.Request.Context() // context.Context を取得
	user, err := ac.AuthService.SignUp(ctx, c)
	if err != nil {
		// エラーレスポンスの処理をここに追加
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ac.AuthService.SendAuthResponse(ctx, c, user, http.StatusCreated)
}

// RefreshAccessTokenHandler リフレッシュトークンを使用してアクセストークンを更新する
func (ac *AuthController) RefreshAccessToken(c *gin.Context) {
    refreshToken := c.GetHeader("Authorization")
    if refreshToken == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not provided"})
        return
    }

    // "Bearer " プレフィックスを削除して refreshToken を抽出
    refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

    // refreshToken を使用して新しい accessToken を取得
    newAccessToken, err := ac.AuthService.RefreshAccessToken(c, refreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // 新しい accessToken をレスポンスとして返す
    c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})
}