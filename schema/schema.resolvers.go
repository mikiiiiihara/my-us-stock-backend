package schema

import (
    "gorm.io/gorm"
    "my-us-stock-backend/schema/generated"
    "my-us-stock-backend/schema/user"
    repoUser "my-us-stock-backend/repository/user"
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
