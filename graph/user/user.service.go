// graph/user/user.service.go

package user

import (
    "context"
    "my-us-stock-backend/graph/user/model"
)

type UserService struct {
    Repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{Repo: repo}
}

// 指定されたIDのユーザーを取得します。
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
    return s.Repo.FindUserByID(ctx, id)
}

// 新規ユーザーを作成します。
func (s *UserService) CreateUser(ctx context.Context, name string, email string) (*model.User, error) {
    return s.Repo.CreateUser(ctx, name, email)
}