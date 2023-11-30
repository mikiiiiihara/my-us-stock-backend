package user

import (
    "context"
    "my-us-stock-backend/schema/generated"
    "my-us-stock-backend/repository/user"
    "my-us-stock-backend/repository/user/model"
    "strconv"
)

// UserService インターフェースの定義
type UserService interface {
    GetUserByID(ctx context.Context, id uint) (*generated.User, error)
    CreateUser(ctx context.Context, name string, email string) (*generated.User, error)
}

type DefaultUserService struct {
    Repo *user.UserRepository
}

func NewUserService(repo *user.UserRepository) UserService {
    return &DefaultUserService{Repo: repo}
}

func (s *DefaultUserService) GetUserByID(ctx context.Context, id uint) (*generated.User, error) {
    modelUser, err := s.Repo.FindUserByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return convertModelUserToGeneratedUser(modelUser), nil
}

func (s *DefaultUserService) CreateUser(ctx context.Context, name string, email string) (*generated.User, error) {
    modelUser, err := s.Repo.CreateUser(ctx, name, email)
    if err != nil {
        return nil, err
    }
    return convertModelUserToGeneratedUser(modelUser), nil
}

// convertModelUserToGeneratedUser は model.User を generated.User に変換します。
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
