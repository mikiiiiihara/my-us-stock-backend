package user

import (
	"context"
	"my-us-stock-backend/app/repository/user/dto"
	"my-us-stock-backend/app/repository/user/model"

	"gorm.io/gorm"
)

// UserRepository インターフェースの定義
type UserRepository interface {
    FindUserByID(ctx context.Context, id uint) (*model.User, error)
    CreateUser(ctx context.Context, dto dto.CreateUserDto) (*model.User, error)
}

// DefaultUserRepository 構造体の定義
type DefaultUserRepository struct {
    DB *gorm.DB
}

// NewUserRepository は DefaultUserRepository の新しいインスタンスを作成します
func NewUserRepository(db *gorm.DB) UserRepository {
    return &DefaultUserRepository{DB: db}
}

// FindUserByID はユーザーをIDによって検索します
func (r *DefaultUserRepository) FindUserByID(ctx context.Context, id uint) (*model.User, error) {
    user := new(model.User)
    result := r.DB.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

// CreateUser は新しいユーザーをデータベースに保存します
func (r *DefaultUserRepository) CreateUser(ctx context.Context, dto dto.CreateUserDto) (*model.User, error) {
    user := &model.User{Name: dto.Name, Email: dto.Email}
    if err := r.DB.Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}