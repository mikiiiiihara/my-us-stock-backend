package user

import (
	"context"
	"my-us-stock-backend/src/schema/generated"
	"strconv"
)

type UserResolver struct {
    UserService UserService
}

func NewResolver(userService UserService) *UserResolver {
    return &UserResolver{UserService: userService}
}

func (r *UserResolver) User(ctx context.Context, idStr string) (*generated.User, error) {
    // string型のIDをuint型に変換
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        // ID変換エラーのハンドリング
        return nil, err
    }

    // GetUserByID に uint 型の ID を渡す
    userModel, err := r.UserService.GetUserByID(ctx, uint(id))
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
func (r *UserResolver) CreateUser(ctx context.Context, name string, email string) (*generated.User, error) {
    userModel, err := r.UserService.CreateUser(ctx, name, email)
    if err != nil {
        return nil, err
    }

    return &generated.User{
        ID:    userModel.ID,  // string型に変換したIDを使用
        Name:  userModel.Name,
        Email: userModel.Email,
    }, nil
}
