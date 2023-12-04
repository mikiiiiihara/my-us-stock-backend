package schema

import (
	repoUser "my-us-stock-backend/src/repository/user"
	"my-us-stock-backend/src/schema/generated"
	"my-us-stock-backend/src/schema/user"

	"gorm.io/gorm"
)

type Resolver struct {
    userResolver *user.Resolver
}

func NewResolver(db *gorm.DB) *Resolver {
    userRepository := repoUser.NewUserRepository(db)
    userService := user.NewUserService(userRepository)
    userResolver := user.NewResolver(userService)

    return &Resolver{
        userResolver: userResolver,
    }
}

func (r *Resolver) Query() generated.QueryResolver {
    return r.userResolver
}

func (r *Resolver) Mutation() generated.MutationResolver {
    return r.userResolver
}
