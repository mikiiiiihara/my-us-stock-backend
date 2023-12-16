package user

import (
	"context"
	"my-us-stock-backend/app/repository/user/dto"
	"my-us-stock-backend/app/repository/user/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テスト用のデータベース設定
func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // テスト用のテーブルを準備
    db.AutoMigrate(&model.User{})

    return db
}

func TestFindUserByID(t *testing.T) {
    db := setupTestDB()
    repo := NewUserRepository(db)

    // テスト用のユーザーを作成
    user := model.User{Name: "Test User", Email: "test@example.com"}
    db.Create(&user)

    // ユーザーをIDで検索
    found, err := repo.FindUserByID(context.Background(), user.ID)
    assert.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, user.Name, found.Name)
    assert.Equal(t, user.Email, found.Email)
}

func TestCreateUser(t *testing.T) {
    db := setupTestDB()
    repo := NewUserRepository(db)

    // 新しいユーザーを作成
        // 新しいデータを作成
	createDto := dto.CreateUserDto{
        Name:   "New User",
        Email: "newuser@example.com",
    }
    created, err := repo.CreateUser(context.Background(), createDto)
    assert.NoError(t, err)
    assert.NotNil(t, created)
    assert.Equal(t, createDto.Name, created.Name)
    assert.Equal(t, createDto.Email, created.Email)

    // データベースでユーザーを確認
    var user model.User
    db.First(&user, created.ID)
    assert.Equal(t, createDto.Name, created.Name)
    assert.Equal(t, createDto.Email, created.Email)
}
