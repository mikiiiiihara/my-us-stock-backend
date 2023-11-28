// graph/user/user.resolver.go

package user

import (
    "context"
    "my-us-stock-backend/graph/user/model"
)

type Resolver struct {
    UserService *UserService
}

func NewResolver(userService *UserService) *Resolver {
    return &Resolver{UserService: userService}
}

func (r *Resolver) User(ctx context.Context, id uint) (*model.User, error) {
    return r.UserService.GetUserByID(ctx, id)
}

// MutationのCreateUserフィールドのResolverです。
func (r *Resolver) CreateUser(ctx context.Context, name string, email string) (*model.User, error) {
    return r.UserService.CreateUser(ctx, name, email)
}
