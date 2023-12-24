package user

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/repository/user/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テスト用のデータベース設定
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // テスト終了時にテーブルの内容をクリアする
    t.Cleanup(func() {
        db.Exec("DELETE FROM users")
    })

    // テスト用のテーブルを準備
    db.AutoMigrate(&model.User{})
    return db
}

func TestFindUserByID(t *testing.T) {
    db := setupTestDB(t)
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
    db := setupTestDB(t)
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

func TestGetUserByEmail(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    // テスト用のユーザーを作成
    user := model.User{Name: "Test User", Email: "test@example.com", Password: "test"}
    db.Create(&user)

    // メールアドレスでユーザーを検索
    found, err := repo.GetUserByEmail(context.Background(), "test@example.com")
    assert.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, user.Name, found.Name)
    assert.Equal(t, user.Email, found.Email)
}

func TestGetAllUserByEmail(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    // 異なるメールアドレスを持つユーザーを作成
    users := []model.User{
        {Name: "Test User 1", Email: "test1@example.com", Password: "test"},
        {Name: "Test User 2", Email: "test2@example.com", Password: "test"},
        {Name: "Test User 3", Email: "test3@example.com", Password: "test"},
    }
    for _, user := range users {
        if err := db.Create(&user).Error; err != nil {
            t.Fatalf("Failed to insert user: %v", err)
        }
    }

    // 特定のメールアドレスでユーザーを検索
    foundUsers, err := repo.GetAllUserByEmail(context.Background(), "test2@example.com")
    assert.NoError(t, err)
    assert.Len(t, foundUsers, 1) // 1人のユーザーが見つかることを期待
    assert.Equal(t, "Test User 2", foundUsers[0].Name)
    assert.Equal(t, "test2@example.com", foundUsers[0].Email)
}
