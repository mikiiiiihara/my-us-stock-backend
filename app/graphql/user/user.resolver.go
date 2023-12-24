package user

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
)

type Resolver struct {
    UserService UserService
}

func NewResolver(userService UserService) *Resolver {
    return &Resolver{UserService: userService}
}

func (r *Resolver) User(ctx context.Context) (*generated.User, error) {
    // GetUserByID に uint 型の ID を渡す
    userModel, err := r.UserService.GetUserByID(ctx)
    if err != nil {
        return nil, err
    }

    return &generated.User{
        ID:    userModel.ID,
        Name:  userModel.Name,
        Email: userModel.Email,
    }, nil
}

// MutationのCreateUserフィールドのResolverです。
func (r *Resolver) CreateUser(ctx context.Context, input generated.CreateUserInput) (*generated.User, error) {
    userModel, err := r.UserService.CreateUser(ctx, input)
    if err != nil {
        return nil, err
    }

    return &generated.User{
        ID:    userModel.ID,  // string型に変換したIDを使用
        Name:  userModel.Name,
        Email: userModel.Email,
    }, nil
}
