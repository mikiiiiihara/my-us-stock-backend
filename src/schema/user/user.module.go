package user

import (
	repoUser "my-us-stock-backend/src/repository/user"

	"gorm.io/gorm"
)

type UserModule struct {
	UserResolver *Resolver
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepoModule := repoUser.NewUserRepositoryModule(db)
	userService := NewUserService(userRepoModule.Repository)
	userResolver := NewResolver(userService)

	return &UserModule{
		UserResolver: userResolver,
	}
}

func (um *UserModule) Query() *Resolver {
	return um.UserResolver
}

func (um *UserModule) Mutation() *Resolver {
	return um.UserResolver
}