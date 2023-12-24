package user

import (
	"encoding/json"
	"my-us-stock-backend/app/repository/user/model"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserE2E(t *testing.T) {
    db := test.SetupTestDB(t)
    router := graphql.SetupGraphQLServer(db, nil)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.User{Name: "Test User", Email: "test@example.com", Password: "abc123"})

    // ダミーのアクセストークンを生成
    token, err := graphql.GenerateTestAccessTokenForUserId(t, 1)
    if err != nil {
        t.Fatalf("Failed to generate test access token: %v", err)
    }

    // GraphQLリクエストの実行
    query := `{ user{ id name email } }`
    w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

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
    err = json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v", err)
    }

    // レスポンスボディの内容の検証
    assert.Equal(t, "1", response.Data.User.ID)
    assert.Equal(t, "Test User", response.Data.User.Name)
    assert.Equal(t, "test@example.com", response.Data.User.Email)
}

func TestCreateUserE2E(t *testing.T) {
	db := test.SetupTestDB(t)
    graphqlServer := graphql.SetupGraphQLServer(db,nil)

	// GraphQLミューテーションの実行
	mutation := `
		mutation {
			createUser(input: {name: "Jane Doe", email: "jane@example.com"}) {
				id
				name
				email
			}
			}	  
	`
	w := graphql.ExecuteGraphQLRequest(graphqlServer, mutation)

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
