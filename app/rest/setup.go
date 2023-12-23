package rest

import (
	repoUser "my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/rest/auth"
	"my-us-stock-backend/app/rest/auth/logic"
	"my-us-stock-backend/app/rest/auth/validation"
	"my-us-stock-backend/app/rest/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupREST は REST API のルートとコントローラを設定します
func SetupREST(r *gin.Engine, db *gorm.DB) {
    // ユーザーリポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)

    // 認証機能
    authLogic := logic.NewAuthLogic()
    userLogic := logic.NewUserLogic()
    responseLogic := logic.NewResponseLogic()
    jwtLogic := logic.NewJWTLogic()
    authValidation := validation.NewAuthValidation()

    // RESTサービス、コントローラの初期化
    userRestService := user.NewUserService(userRepo)
    userController := user.NewUserController(userRestService)

    // 認証サービスの初期化
    authService := auth.NewAuthService(userRepo, authLogic, userLogic, responseLogic, jwtLogic, authValidation)
    authController := auth.NewAuthController(authService)

    // RESTコントローラのルートを設定
    r.GET("/api/users/:id", userController.GetUser)
    r.POST("/api/users", userController.CreateUser)
    // 認証用
    r.POST("/api/v1/signin", authController.SignIn)
    r.POST("/api/v1/signup", authController.SignUp)
}
