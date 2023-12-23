package logic_test

import (
	"fmt"
	"my-us-stock-backend/app/rest/auth/logic"
	"net/http"
	"testing"
)

func TestGetUserIdFromContextNotAuthenticationTokenError(t *testing.T) {
	expectedUserId := 1
	invalidToken := ""

	// リクエストの生成
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/todo", nil)

	req.Header.Add("Authorization", invalidToken)

	// テスト実行
	authLogic := logic.NewAuthLogic()
	actual, err := authLogic.GetUserIdFromContext(req)

	expectedError := "not token"

	if err.Error() != expectedError || actual != 0 {
		t.Errorf("actual %v\nwant %v", actual, expectedUserId)
	}
}

func TestGetUserIdFromContextEmptyTokenError(t *testing.T) {
	expectedUserId := 1
	invalidToken := ""

	// リクエストの生成
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/todo", nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidToken))

	// テスト実行
	authLogic := logic.NewAuthLogic()
	actual, err := authLogic.GetUserIdFromContext(req)

	expectedError := "トークンが空文字です。"

	if err.Error() != expectedError || actual != 0 {
		fmt.Print(err)
		t.Errorf("actual %v\nwant %v", actual, expectedUserId)
	}
}
