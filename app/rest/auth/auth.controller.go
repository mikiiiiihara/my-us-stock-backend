package auth

import (
	"net/http"

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

