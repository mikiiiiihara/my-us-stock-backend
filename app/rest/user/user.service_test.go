package user

import (
	"context"
	userModel "my-us-stock-backend/app/database/model"
	repoUser "my-us-stock-backend/app/repository/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository は UserRepository のモックです。
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, id uint) (*userModel.User, error) {
    args := m.Called(ctx, id)
    // 戻り値の型が *userModel.User であることを確認
    return args.Get(0).(*userModel.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, createDto repoUser.CreateUserDto) (*userModel.User, error) {
    args := m.Called(ctx, createDto)
    // 戻り値の型が *userModel.User であることを確認
    return args.Get(0).(*userModel.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
    args := m.Called(ctx, email)
    // 戻り値の型が *userModel.User であることを確認
    var user *userModel.User
    if args.Get(0) != nil {
        user = args.Get(0).(*userModel.User)
    }
    return user, args.Error(1)
}

func (m *MockUserRepository) GetAllUserByEmail(ctx context.Context, email string) ([]*userModel.User, error) {
    args := m.Called(ctx, email)
    // 戻り値の型が []*userModel.User であることを確認
    return args.Get(0).([]*userModel.User), args.Error(1)
}


var _ repoUser.UserRepository = (*MockUserRepository)(nil)

// TestGetUserByID は GetUserByID メソッドのテストです。
func TestGetUserByID(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)   // repoUser エイリアスを使用

	mockUser := &userModel.User{
		Model: gorm.Model{ID: 1},  // gorm.Model で ID を設定
		Name:  "John Doe",
		Email: "john@example.com",
	}
    mockRepo.On("FindUserByID", mock.Anything, uint(1)).Return(mockUser, nil)

    result, err := service.GetUserByID(context.Background(), 1)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "John Doe", result.Name)
    assert.Equal(t, "john@example.com", result.Email)
}

// TestCreateUserService は CreateUser メソッドのテストです。
func TestCreateUserService(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)  // repoUser エイリアスを使用
    createUserDto := repoUser.CreateUserDto{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }
	mockUser := &userModel.User{
		Model: gorm.Model{ID: 1},  // gorm.Model で ID を設定
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
    mockRepo.On("CreateUser", mock.Anything, createUserDto).Return(mockUser, nil)

    createUserInput := CreateUserInput{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }

    result, err := service.CreateUser(context.Background(), createUserInput)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, createUserInput.Name, result.Name)
    assert.Equal(t, createUserInput.Email, result.Email)
}
