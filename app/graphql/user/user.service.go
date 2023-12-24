package user

import (
	"context"
	"fmt"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
	"my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/repository/user/dto"
	"my-us-stock-backend/app/repository/user/model"
	"strconv"
)

// UserService インターフェースの定義
type UserService interface {
    GetUserByID(ctx context.Context) (*generated.User, error)
    CreateUser(ctx context.Context, input generated.CreateUserInput) (*generated.User, error)
}

// DefaultUserService 構造体の定義
type DefaultUserService struct {
    Repo user.UserRepository // インターフェースを利用
    Auth auth.AuthService    // 認証サービスのインターフェース
}

// NewUserService は DefaultUserService の新しいインスタンスを作成します
func NewUserService(repo user.UserRepository, auth auth.AuthService) UserService {
    return &DefaultUserService{Repo: repo, Auth: auth}
}

// GetUserByID はユーザーをIDによって検索します
func (s *DefaultUserService) GetUserByID(ctx context.Context) (*generated.User, error) {
    // cookieからアクセストークンを取得
    accessToken, _ := ctx.Value(utils.CookieKey).(string)
    // アクセストークンの検証
    userId, _ := s.Auth.FetchUserIdAccessToken(accessToken)
    modelUser, err := s.Repo.FindUserByID(ctx, userId)
    if err != nil {
        return nil, err
    }
    fmt.Println("------------")
    fmt.Println(userId)
    fmt.Println("------------")
    fmt.Println(modelUser)
    return convertModelUserToGeneratedUser(modelUser), nil
}

// CreateUser は新しいユーザーを作成します
func (s *DefaultUserService) CreateUser(ctx context.Context, input generated.CreateUserInput) (*generated.User, error) {
    // 更新用DTOの作成
    createDto := dto.CreateUserDto{
        Name: input.Name,
        Email: input.Email,
    }
    modelUser, err := s.Repo.CreateUser(ctx, createDto)
    if err != nil {
        return nil, err
    }
    return convertModelUserToGeneratedUser(modelUser), nil
}

// convertModelUserToGeneratedUser は model.User を generated.User に変換します
func convertModelUserToGeneratedUser(modelUser *model.User) *generated.User {
    if modelUser == nil {
        return nil
    }
    return &generated.User{
        ID:    strconv.FormatUint(uint64(modelUser.ID), 10),
        Name:  modelUser.Name,
        Email: modelUser.Email,
    }
}
