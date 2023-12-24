package logic

import (
	"my-us-stock-backend/app/database/model"
	"os"
	"testing"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessToken(t *testing.T) {
    // 環境変数の設定
    os.Setenv("JWT_KEY", "testkey")

    jwtLogic := NewJWTLogic()

    user := &model.User{
        Name: "Test User",
		Email: "test@test.com",
    }

    tokenString, err := jwtLogic.CreateAccessToken(user)
    assert.NoError(t, err)
    assert.NotEmpty(t, tokenString)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("testkey"), nil
    })

    // ここでのエラーチェック
    if assert.NoError(t, err) {
        if assert.NotNil(t, token) {
            claims, ok := token.Claims.(jwt.MapClaims)
            assert.True(t, ok)
            assert.Equal(t, true, claims["admin"])
            assert.Equal(t, user.Name, claims["name"])
            assert.Equal(t, user.ID, uint(claims["id"].(float64)))

            expiration := time.Unix(int64(claims["exp"].(float64)), 0)
            assert.WithinDuration(t, time.Now().Add(time.Minute*15), expiration, time.Second)
        }
    }
}

func TestCreateRefreshToken(t *testing.T) {
    // 環境変数の設定
    os.Setenv("JWT_KEY", "testkey")

    jwtLogic := NewJWTLogic()

    user := &model.User{
        Name: "Test User",
		Email: "test@test.com",
    }

    tokenString, err := jwtLogic.CreateRefreshToken(user)
    assert.NoError(t, err)
    assert.NotEmpty(t, tokenString)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("testkey"), nil
    })

    // ここでのエラーチェック
    if assert.NoError(t, err) {
        if assert.NotNil(t, token) {
            claims, ok := token.Claims.(jwt.MapClaims)
            assert.True(t, ok)
            assert.Equal(t, user.ID, uint(claims["id"].(float64)))

            expiration := time.Unix(int64(claims["exp"].(float64)), 0)
            assert.WithinDuration(t, time.Now().Add(time.Hour * 24 * 7), expiration, time.Second)
        }
    }
}
