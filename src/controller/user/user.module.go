package user

import (
	repoUser "my-us-stock-backend/src/repository/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserModule struct {
    UserController *UserController
}

func NewUserModule(db *gorm.DB) *UserModule {
    userRepository := repoUser.NewUserRepository(db)
    userService := NewUserService(userRepository)
    userController := NewUserController(userService)

    return &UserModule{
        UserController: userController,
    }
}

func (um *UserModule) RegisterRoutes(router *gin.Engine) {
    router.GET("/api/users/:id", um.UserController.GetUser)
    router.POST("/api/users", um.UserController.CreateUser)
}
