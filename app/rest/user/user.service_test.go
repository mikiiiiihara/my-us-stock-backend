package user

import (
	"context"
	repoUser "my-us-stock-backend/app/repository/user"
	userModel "my-us-stock-backend/app/repository/user/model"
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

func (m *MockUserRepository) CreateUser(ctx context.Context, name string, email string) (*userModel.User, error) {
    args := m.Called(ctx, name, email)
    // 戻り値の型が *userModel.User であることを確認
    return args.Get(0).(*userModel.User), args.Error(1)
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

	mockUser := &userModel.User{
		Model: gorm.Model{ID: 1},  // gorm.Model で ID を設定
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
    mockRepo.On("CreateUser", mock.Anything, "Jane Doe", "jane@example.com").Return(mockUser, nil)

    result, err := service.CreateUser(context.Background(), "Jane Doe", "jane@example.com")
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Jane Doe", result.Name)
    assert.Equal(t, "jane@example.com", result.Email)
}
