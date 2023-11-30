package user

import (
    "context"
    "gorm.io/gorm"
    "my-us-stock-backend/repository/user/model"
)

type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) FindUserByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    result := r.DB.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

// Create は新しいユーザーをデータベースに保存します。
func (r *UserRepository) CreateUser(ctx context.Context, name string, email string) (*model.User, error) {
    user := &model.User{Name: name, Email: email}
    if err := r.DB.Create(user).Error; err != nil {
        return nil, err
    }
    return user, nil
}