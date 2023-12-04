package user_test

import (
	"bytes"
	"encoding/json"
	"my-us-stock-backend/src/controller/user"
	RepoUser "my-us-stock-backend/src/repository/user"
	"my-us-stock-backend/src/repository/user/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

// テスト用のユーザーコントローラをセットアップ
func setupUserController(db *gorm.DB) *user.UserController {
    userRepository := RepoUser.NewUserRepository(db)
    userService := user.NewUserService(userRepository)
    return user.NewUserController(userService)
}

// GetUserのe2eテスト
func TestGetUserE2E(t *testing.T) {
    db := setupTestDB()
    controller := setupUserController(db)

    // テスト用のユーザーを作成
    user := model.User{Name: "Test User", Email: "test@example.com"}
    db.Create(&user)

    // Ginのルーターをセットアップ
    router := gin.Default()
    router.GET("/api/users/:id", controller.GetUser)

    // テストリクエストの作成
    req, _ := http.NewRequest("GET", "/api/users/1", nil)
    w := httptest.NewRecorder()

    // リクエストの実行
    router.ServeHTTP(w, req)

    // レスポンスの検証
    assert.Equal(t, http.StatusOK, w.Code)
    // レスポンスボディの解析
    var response struct {
        ID    string `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v", err)
    }

    // レスポンスボディの内容の検証
    assert.Equal(t, "1", response.ID)
    assert.Equal(t, "Test User", response.Name)
    assert.Equal(t, "test@example.com", response.Email)
}

// CreateUserのe2eテスト
func TestCreateUserE2E(t *testing.T) {
    db := setupTestDB()
    controller := setupUserController(db)

    // テスト用のユーザーを作成
    seedUser := model.User{Name: "Test User", Email: "test@example.com"}
    db.Create(&seedUser)

    router := gin.Default()
    router.POST("/api/users", controller.CreateUser)

    // リクエストボディの作成
    body, _ := json.Marshal(user.CreateUserInput{
        Name: "Jane Doe",
        Email: "janedoe@example.com",
    })

    req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    // レスポンスボディの解析
    var newUser struct {
        ID    string `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    err := json.Unmarshal(w.Body.Bytes(), &newUser)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v", err)
    }

    // レスポンスボディの内容の検証
    assert.NotEmpty(t, newUser.ID)
    assert.Equal(t, "Jane Doe", newUser.Name)
    assert.Equal(t, "janedoe@example.com", newUser.Email)
}
