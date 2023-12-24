package graphql

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

// ExecuteGraphQLRequest はGraphQLリクエストを実行し、レスポンスを返す
func ExecuteGraphQLRequest(h http.Handler, query string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
	})
	req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

// テスト用のダミーアクセストークンを用いてGraphQLリクエストを実行する
func ExecuteGraphQLRequestWithToken(url, query, token string) *httptest.ResponseRecorder {
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

    w := httptest.NewRecorder()
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
    return w
}

// 任意のuserIdを含むアクセストークンを作成する
func GenerateTestAccessTokenForUserId(t *testing.T,userId uint) (string, error) {
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
