package controller

import (
	"my-us-stock-backend/src/controller/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControllerModule struct {
    UserModule *user.UserModule
}

// NewControllerModule は新しいControllerModuleインスタンスを作成します
func NewControllerModule(db *gorm.DB) *ControllerModule {
    userModule := user.NewUserModule(db)
    // 他のモジュールの初期化

    return &ControllerModule{
        UserModule: userModule,
        // 他のモジュールのインスタンス化
    }
}

// RegisterRoutes はGinルーターに対してコントローラのルートを登録します
func (cr *ControllerModule) RegisterRoutes(router *gin.Engine) {
    cr.UserModule.RegisterRoutes(router)
    // 他のモジュールのルート登録
}