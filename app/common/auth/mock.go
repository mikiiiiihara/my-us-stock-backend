package auth

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockAuthService は AuthService のモックです。
type MockAuthService struct {
	mock.Mock
}

// NewMockAuthService は新しい NewMockAuthService を作成し、初期設定を行います。
func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (m *MockAuthService) FetchUserIdAccessToken(ctx context.Context) (uint, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuthService) RefreshAccessToken(c *gin.Context) (string, error) {
    args := m.Called(c)
    return args.String(0), args.Error(1)
}

func (m *MockAuthService) SignIn(ctx context.Context, c *gin.Context) (*model.User, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) SignUp(ctx context.Context, c *gin.Context) (*model.User, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) SendAuthResponse(ctx context.Context, c *gin.Context, user *model.User, code int) {
    m.Called(ctx, c, user, code)
}