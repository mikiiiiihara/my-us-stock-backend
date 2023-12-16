package user

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	repoUser "my-us-stock-backend/app/repository/user"
	"my-us-stock-backend/app/repository/user/dto"
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

// CreateUser のモック関数を修正
func (m *MockUserRepository) CreateUser(ctx context.Context, createDto dto.CreateUserDto) (*userModel.User, error) {
    args := m.Called(ctx, createDto)
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
    service := NewUserService(mockRepo)

    createUserInput := generated.CreateUserInput{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }

    mockUser := &userModel.User{
        Model: gorm.Model{ID: 1},
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }

    mockRepo.On("CreateUser", mock.Anything, dto.CreateUserDto{Name: "Jane Doe", Email: "jane@example.com"}).Return(mockUser, nil)

    result, err := service.CreateUser(context.Background(), createUserInput)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Jane Doe", result.Name)
    assert.Equal(t, "jane@example.com", result.Email)
}
