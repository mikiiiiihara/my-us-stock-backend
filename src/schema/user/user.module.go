package user

import (
	repoUser "my-us-stock-backend/src/repository/user"
	"my-us-stock-backend/src/schema/generated"

	"gorm.io/gorm"
)

type UserModule struct {
	Resolver *Resolver
}

func NewUserModule(db *gorm.DB) *UserModule {
	userRepository := repoUser.NewUserRepository(db)
	userService := NewUserService(userRepository)
	userResolver := NewResolver(userService)

	return &UserModule{
		Resolver: userResolver,
	}
}

func (um *UserModule) Query() generated.QueryResolver {
	return um.Resolver
}

func (um *UserModule) Mutation() generated.MutationResolver {
	return um.Resolver
}
