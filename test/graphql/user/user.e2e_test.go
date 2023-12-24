package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"my-us-stock-backend/app/repository/user/model"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/stretchr/testify/assert"
)

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

func executeGraphQLRequestWithToken(url, query, token string) *httptest.ResponseRecorder {
    requestBody, _ := json.Marshal(map[string]interface{}{
        "query": query,
    })

    req, _ := http.NewRequest("POST", url+"/graphql", bytes.NewBuffer(requestBody))
    req.Header.Set("Content-Type", "application/json")
    req.AddCookie(&http.Cookie{Name: "access_token", Value: token})

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Printf("Error making request: %v", err)
        return nil
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return nil
    }
    log.Printf("Response Status: %s", resp.Status)
    log.Printf("Response Body: %s", string(body))

    w := httptest.NewRecorder()
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
    return w
}



func generateTestAccessTokenForUserId(t *testing.T,userId uint) (string, error) {
    // JWTの構造を準備
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)

    // ダミーのユーザー情報をクレームにセット
    claims["admin"] = true
    claims["sub"] = strconv.Itoa(int(userId)) + "test@example.com" + "Test User"
    claims["id"] = userId
    claims["name"] = "Test User"
    claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

    // 署名してトークンを生成
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func TestUserE2E(t *testing.T) {
    db := test.SetupTestDB(t)
    router := graphql.SetupGraphQLServer(db, nil)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.User{Name: "Test User", Email: "test@example.com", Password: "abc123"})

    // ダミーのアクセストークンを生成
    token, err := generateTestAccessTokenForUserId(t, 1)
    if err != nil {
        t.Fatalf("Failed to generate test access token: %v", err)
    }

    // GraphQLリクエストの実行
    query := `{ user{ id name email } }`
    w := executeGraphQLRequestWithToken(ts.URL, query, token)

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
	fmt.Println("------------")
    fmt.Println(response)
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
