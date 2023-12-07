package user

import (
	"bytes"
	"context"
	"encoding/json"
	"my-us-stock-backend/src/repository/user/model"
	"my-us-stock-backend/src/schema/currency"
	"my-us-stock-backend/src/schema/generated"
	"my-us-stock-backend/src/schema/user"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CustomQueryResolver struct {
    UserModule     *user.UserModule
    CurrencyModule *currency.CurrencyModule
}

func NewCustomQueryResolver(userModule *user.UserModule, currencyModule *currency.CurrencyModule) *CustomQueryResolver {
    return &CustomQueryResolver{
        UserModule:     userModule,
        CurrencyModule: currencyModule,
    }
}

func (r *CustomQueryResolver) Query() generated.QueryResolver {
    return r
}

func (r *CustomQueryResolver) Mutation() generated.MutationResolver {
    return r.UserModule.Mutation()
}

// Userメソッドの実装
func (r *CustomQueryResolver) User(ctx context.Context, id string) (*generated.User, error) {
    return r.UserModule.Query().User(ctx, id)
}

// GetCurrentUsdJpyメソッドの実装
func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
    return r.CurrencyModule.Query().GetCurrentUsdJpy(ctx)
}


// setupTestDB はテスト用のデータベースをセットアップします
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&model.User{})
	return db
}

// setupGraphQLServer はテスト用のGraphQLサーバーをセットアップします
func setupGraphQLServer(db *gorm.DB) *handler.Server {
	// テスト用の環境変数を設定
	os.Setenv("CURRENCY_URL", "http://test.url")
		
    userModule := user.NewUserModule(db)
    currencyModule := currency.NewCurrencyModule()

    customResolver := NewCustomQueryResolver(userModule, currencyModule)

    return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: customResolver}))
}

// executeGraphQLRequest はGraphQLリクエストを実行し、レスポンスを返します
func executeGraphQLRequest(h http.Handler, query string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
	})
	req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestUserE2E(t *testing.T) {
	db := setupTestDB(t)
	graphqlServer := setupGraphQLServer(db)

	// テスト用のユーザーを作成
	db.Create(&model.User{Name: "Test User", Email: "test@example.com"})

	// GraphQLリクエストの実行
	query := `{ user(id: "1") { id name email } }`
	w := executeGraphQLRequest(graphqlServer, query)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			User struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"user"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.Equal(t, "1", response.Data.User.ID)
	assert.Equal(t, "Test User", response.Data.User.Name)
	assert.Equal(t, "test@example.com", response.Data.User.Email)
}

func TestCreateUserE2E(t *testing.T) {
	db := setupTestDB(t)
	graphqlServer := setupGraphQLServer(db)

	// GraphQLミューテーションの実行
	mutation := `
		mutation {
			createUser(name: "Jane Doe", email: "jane@example.com") {
				id
				name
				email
			}
		}
	`
	w := executeGraphQLRequest(graphqlServer, mutation)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			CreateUser struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"createUser"`
		} `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.NotEmpty(t, response.Data.CreateUser.ID)
	assert.Equal(t, "Jane Doe", response.Data.CreateUser.Name)
	assert.Equal(t, "jane@example.com", response.Data.CreateUser.Email)
}
