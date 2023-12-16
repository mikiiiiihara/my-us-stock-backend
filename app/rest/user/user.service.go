package user

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/repository/user/dto"
	"my-us-stock-backend/app/repository/user/model"
	"my-us-stock-backend/app/rest/user/input"
	"strconv"
)

// UserService インターフェースの定義
type UserService interface {
    GetUserByID(ctx context.Context, id uint) (*generated.User, error)
    CreateUser(ctx context.Context, input input.CreateUserInput) (*generated.User, error)
}

// DefaultUserService 構造体の定義
type DefaultUserService struct {
    Repo user.UserRepository // インターフェースを利用
}

// NewUserService は DefaultUserService の新しいインスタンスを作成します
func NewUserService(repo user.UserRepository) UserService {
    return &DefaultUserService{Repo: repo}
}

// GetUserByID はユーザーをIDによって検索します
func (s *DefaultUserService) GetUserByID(ctx context.Context, id uint) (*generated.User, error) {
    modelUser, err := s.Repo.FindUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return convertModelUserToGeneratedUser(modelUser), nil
}

// CreateUser は新しいユーザーを作成します
func (s *DefaultUserService) CreateUser(ctx context.Context, input input.CreateUserInput) (*generated.User, error) {
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
