package user

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService の定義
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*generated.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*generated.User), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, input CreateUserInput) (*generated.User, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.User), args.Error(1)
}

// MockUserService が UserService インターフェースを実装することを確認
var _ UserService = (*MockUserService)(nil)
// ユーザー取得のテスト
func TestUserController_GetUser(t *testing.T) {
    mockService := new(MockUserService)
    controller := NewUserController(mockService)

    // モックの戻り値として期待される generated.User 型のインスタンスを作成
    expectedUser := &generated.User{ID: "1", Name: "John Doe", Email: "johndoe@example.com"}
    mockService.On("GetUserByID", mock.Anything, uint(1)).Return(expectedUser, nil)

    // ここで直接メソッドを呼び出す
    user, err := controller.UserService.GetUserByID(context.Background(), 1)

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockService.AssertExpectations(t)
}

// ユーザー作成のテスト
func TestUserController_CreateUser(t *testing.T) {
    mockService := new(MockUserService)
    controller := NewUserController(mockService)

    expectedUser := &generated.User{Name: "Jane Doe", Email: "janedoe@example.com"}

    createUserInput := CreateUserInput{
        Name:  "Jane Doe",
        Email: "jane@example.com",
    }
    mockService.On("CreateUser", mock.Anything, createUserInput).Return(expectedUser, nil)

    // ここで直接メソッドを呼び出す
    user, err := controller.UserService.CreateUser(context.Background(), createUserInput)

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockService.AssertExpectations(t)
}
