package user

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/graphql/generated"
	"strconv"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

type Resolver struct {
    UserService UserService
    AuthService auth.AuthService
}

func NewResolver(userService UserService) *Resolver {
    return &Resolver{UserService: userService}
}

func (r *Resolver) User(ctx context.Context, idStr string) (*generated.User, error) {
    // GraphQLリクエストのヘッダーからアクセストークンを取得
    opCtx := graphql.GetOperationContext(ctx)
    authorizationValue := opCtx.Headers.Get("Authorization")
    accessToken := strings.TrimPrefix(authorizationValue, "Bearer ")
    fmt.Println(accessToken)
    
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
