package auth

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthController_SignIn(t *testing.T) {
    mockAuthService := auth.NewMockAuthService()
    controller := NewAuthController(mockAuthService)

    mockUser := &model.User{Name: "Test User", Email: "test@example.com"}
    mockAuthService.On("SignIn", mock.Anything, mock.Anything).Return(mockUser, nil)

    // ここで直接メソッドを呼び出す
    user, err := controller.AuthService.SignIn(context.Background(), nil)

    assert.NoError(t, err)
    assert.Equal(t, mockUser, user)
    mockAuthService.AssertExpectations(t)
}

func TestAuthController_SignUp(t *testing.T) {
    mockAuthService := auth.NewMockAuthService()
    controller := NewAuthController(mockAuthService)

    mockUser := &model.User{Name: "New User", Email: "newuser@example.com"}
    mockAuthService.On("SignUp", mock.Anything, mock.Anything).Return(mockUser, nil)

    // ここで直接メソッドを呼び出す
    user, err := controller.AuthService.SignUp(context.Background(), nil)
    assert.NoError(t, err)
    assert.Equal(t, mockUser, user)
    mockAuthService.AssertExpectations(t)
}

func TestAuthController_RefreshAccessToken(t *testing.T) {
    // MockAuthServiceのインスタンスを作成
    mockAuthService := auth.NewMockAuthService()
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
