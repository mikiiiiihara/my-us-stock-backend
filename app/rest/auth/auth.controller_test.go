package auth

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	userModel "my-us-stock-backend/app/repository/user/model"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
    mock.Mock
}

func (m *MockAuthService) GetUserIdFromToken(w http.ResponseWriter, r *http.Request) (int, error) {
    args := m.Called(w, r)
    return args.Int(0), args.Error(1)
}

func (m *MockAuthService) SignIn(ctx context.Context, c *gin.Context) (*userModel.User, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(*userModel.User), args.Error(1)
}

func (m *MockAuthService) SignUp(ctx context.Context, c *gin.Context) (*userModel.User, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(*userModel.User), args.Error(1)
}

func (m *MockAuthService) SendAuthResponse(ctx context.Context, c *gin.Context, user *userModel.User, code int) {
    m.Called(ctx, c, user, code)
}

func (m *MockAuthService) RefreshAccessToken(c *gin.Context) (string, error) {
    args := m.Called(c)
    return args.String(0), args.Error(1)
}

// FetchUserIdAccessTokenのモックメソッド
func (m *MockAuthService) FetchUserIdAccessToken(token string) (uint, error) {
    args := m.Called(token)
    return args.Get(0).(uint), args.Error(1)
}

var _ auth.AuthService = (*MockAuthService)(nil)

func TestAuthController_SignIn(t *testing.T) {
    mockAuthService := new(MockAuthService)
    controller := NewAuthController(mockAuthService)

    mockUser := &userModel.User{Name: "Test User", Email: "test@example.com"}
    mockAuthService.On("SignIn", mock.Anything, mock.Anything).Return(mockUser, nil)

    // ここで直接メソッドを呼び出す
    user, err := controller.AuthService.SignIn(context.Background(), nil)

    assert.NoError(t, err)
    assert.Equal(t, mockUser, user)
    mockAuthService.AssertExpectations(t)
}

func TestAuthController_SignUp(t *testing.T) {
    mockAuthService := new(MockAuthService)
    controller := NewAuthController(mockAuthService)

    mockUser := &userModel.User{Name: "New User", Email: "newuser@example.com"}
    mockAuthService.On("SignUp", mock.Anything, mock.Anything).Return(mockUser, nil)

    // ここで直接メソッドを呼び出す
    user, err := controller.AuthService.SignUp(context.Background(), nil)
    assert.NoError(t, err)
    assert.Equal(t, mockUser, user)
    mockAuthService.AssertExpectations(t)
}

func TestAuthController_RefreshAccessToken(t *testing.T) {
    // MockAuthServiceのインスタンスを作成
    mockAuthService := new(MockAuthService)
    controller := NewAuthController(mockAuthService)

    // モックのリフレッシュトークン
    refreshToken := "mockRefreshToken"

    // モックの新しいアクセストークン
    newAccessToken := "newAccessToken"

    // MockAuthServiceのRefreshAccessTokenメソッドをモック化
    mockAuthService.On("RefreshAccessToken", mock.Anything).Return(newAccessToken, nil)

    // テスト用のgin.Contextを作成
    c := &gin.Context{
        Request: &http.Request{
            Header: http.Header{"Cookie": []string{"refresh_token=" + refreshToken}},
        },
    }

    // モックのRefreshAccessTokenメソッドを呼び出す
    returnedAccessToken, err := controller.AuthService.RefreshAccessToken(c)




    // エラーが発生しないことを検証
    assert.NoError(t, err)

    // 期待される新しいアクセストークンが返されたことを検証
    assert.Equal(t, newAccessToken, returnedAccessToken)

    // MockAuthServiceのExpectationsを検証
    mockAuthService.AssertExpectations(t)
}
