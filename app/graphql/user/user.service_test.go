package user

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	userModel "my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/graphql/utils"
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

// CreateUser のモック関数を修正
func (m *MockUserRepository) CreateUser(ctx context.Context, createDto repoUser.CreateUserDto) (*userModel.User, error) {
    args := m.Called(ctx, createDto)
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
	mockAuth := auth.NewMockAuthService()
    service := NewUserService(mockRepo, mockAuth)

    // モックの期待値設定
    testAccessToken := "testAccessToken"
    expectedUserID := uint(1)  // 明示的に uint 型を使用

    mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(expectedUserID, nil)
    mockUser := &userModel.User{
        Model: gorm.Model{ID: 1},
        Name:  "John Doe",
        Email: "john@example.com",
    }
    mockRepo.On("FindUserByID", mock.Anything, expectedUserID).Return(mockUser, nil)

    // contextにアクセストークンを設定
    ctx := context.WithValue(context.Background(), utils.CookieKey, testAccessToken)

    // テスト対象メソッドの実行
    result, err := service.GetUserByID(ctx)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "John Doe", result.Name)
    assert.Equal(t, "john@example.com", result.Email)

    // モックの呼び出しを検証
    mockRepo.AssertExpectations(t)
    mockAuth.AssertExpectations(t)
}


// TestCreateUserService は CreateUser メソッドのテストです。
func TestCreateUserService(t *testing.T) {
    mockRepo := new(MockUserRepository)
	mockAuth := auth.NewMockAuthService()
    service := NewUserService(mockRepo, mockAuth)

    createUserInput := generated.CreateUserInput{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }

    mockUser := &userModel.User{
        Model: gorm.Model{ID: 1},
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }

    mockRepo.On("CreateUser", mock.Anything, repoUser.CreateUserDto{Name: "Jane Doe", Email: "jane@example.com"}).Return(mockUser, nil)

    result, err := service.CreateUser(context.Background(), createUserInput)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "Jane Doe", result.Name)
    assert.Equal(t, "jane@example.com", result.Email)
}
