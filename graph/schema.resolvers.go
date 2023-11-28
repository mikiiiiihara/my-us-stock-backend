package graph

import (
    "gorm.io/gorm"
    "my-us-stock-backend/graph/generated"
    "my-us-stock-backend/graph/user"
)

type Resolver struct {
    userResolver *user.Resolver
}

func NewResolver(db *gorm.DB) *Resolver {
    userRepository := user.NewUserRepository(db)
    userService := user.NewUserService(userRepository)
    userResolver := user.NewResolver(userService)

    return &Resolver{
        userResolver: userResolver,
    }
}

// Mutation メソッドを追加
func (r *Resolver) Mutation() generated.MutationResolver {
    return r.userResolver
}
