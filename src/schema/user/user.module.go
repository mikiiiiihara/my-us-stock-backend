package user

import (
	repoUser "my-us-stock-backend/src/repository/user"

	"gorm.io/gorm"
)

type UserModule struct {
	UserResolver *UserResolver
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepoModule := repoUser.NewUserRepositoryModule(db)
	userService := NewUserService(userRepoModule.Repository)
	userResolver := NewResolver(userService)

	return &UserModule{
		UserResolver: userResolver,
	}
}

func (um *UserModule) Query() *UserResolver {
	return um.UserResolver
}

func (um *UserModule) Mutation() *UserResolver {
	return um.UserResolver
}