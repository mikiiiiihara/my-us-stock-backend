package rest

import (
	repoUser "my-us-stock-backend/repository/user"
	"my-us-stock-backend/rest/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupREST は REST API のルートとコントローラを設定します
func SetupREST(r *gin.Engine, db *gorm.DB) {
    // ユーザーリポジトリの初期化
    userRepo := repoUser.NewUserRepository(db)

    // RESTサービス、コントローラの初期化
    userRestService := user.NewUserService(userRepo)
    userController := user.NewUserController(userRestService)

    // RESTコントローラのルートを設定
    r.GET("/api/users/:id", userController.GetUser)
    r.POST("/api/users", userController.CreateUser)
}
