package controller

import (
	"my-us-stock-backend/controller/user"
	repoUser "my-us-stock-backend/repository/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControllerRegistry struct {
    UserController *user.UserController
    // 他のコントローラもここに追加
}

// NewControllerRegistry は新しいControllerRegistryインスタンスを作成します
func NewControllerRegistry(db *gorm.DB) *ControllerRegistry {
    // リポジトリのインスタンス化
    userRepository := repoUser.NewUserRepository(db)

    // サービスのインスタンス化
    userService := user.NewUserService(userRepository)

    // コントローラのインスタンス化
    userController := user.NewUserController(userService)

    return &ControllerRegistry{
        UserController: userController,
        // 他のコントローラのインスタンス化
    }
}

// RegisterRoutes はGinルーターに対してコントローラのルートを登録します
func (cr *ControllerRegistry) RegisterRoutes(router *gin.Engine) {
    // UserControllerのルートを設定
    router.GET("/api/users/:id", cr.UserController.GetUser)
    router.POST("/api/users", cr.UserController.CreateUser)
    // 他のコントローラのルートも同様に登録
}
