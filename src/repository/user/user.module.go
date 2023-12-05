package user

import (
	"gorm.io/gorm"
)

type UserRepositoryModule struct {
	Repository UserRepository
}

func NewUserRepositoryModule(db *gorm.DB) *UserRepositoryModule {
	userRepository := NewUserRepository(db)
	return &UserRepositoryModule{
		Repository: userRepository,
	}
}
