package user

import (
	repoUser "my-us-stock-backend/repository/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserModule struct {
    UserController *UserController
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepoModule := repoUser.NewUserRepositoryModule(db)
	userService := NewUserService(userRepoModule.Repository)
    userController := NewUserController(userService)

    return &UserModule{
        UserController: userController,
    }
}

func (um *UserModule) RegisterRoutes(router *gin.Engine) {
    router.GET("/api/users/:id", um.UserController.GetUser)
    router.POST("/api/users", um.UserController.CreateUser)
}
