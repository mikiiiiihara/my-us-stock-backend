package user

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"gorm.io/gorm"
)

// UserRepository インターフェースの定義
type UserRepository interface {
    FindUserByID(ctx context.Context, id uint) (*model.User, error)
    CreateUser(ctx context.Context, dto CreateUserDto) (*model.User, error)
    GetUserByEmail(ctx context.Context, email string) (*model.User, error)
    GetAllUserByEmail(ctx context.Context, email string) ([]*model.User, error)
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
func (r *DefaultUserRepository) CreateUser(ctx context.Context, dto CreateUserDto) (*model.User, error) {
    user := &model.User{Name: dto.Name, Email: dto.Email, Password: dto.Password}
    if err := r.DB.Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}

// GetUserByEmail はemailに紐づくユーザーを取得します
func (r *DefaultUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    user := new(model.User)
    result := r.DB.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

// GetAllUserByEmail はemailに紐づく全てのユーザーを取得します
func (r *DefaultUserRepository) GetAllUserByEmail(ctx context.Context, email string) ([]*model.User, error) {
    var users []*model.User
    result := r.DB.Where("email = ?", email).Find(&users)
    if result.Error != nil {
        return nil, result.Error
    }
    return users, nil
}