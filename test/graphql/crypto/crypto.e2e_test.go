package crypto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-us-stock-backend/app/database/model"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price/crypto"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockHTTPTransport は http.RoundTripper のインターフェースを満たすモック実装です。
type MockHTTPTransport struct {
    RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip は http.RoundTripper の RoundTrip メソッドを模倣します。
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.RoundTripFunc(req)
}

func TestCryptosE2E(t *testing.T) {
	db := test.SetupTestDB()

	// モックの HTTP レスポンスを設定
	mockCryptoPrice := `{
		"success": 1,
		"data": {
			"sell": "5958001",
			"buy": "5958000",
			"open": "6052000",
			"high": "6127930",
			"low": "5900000",
			"last": "5956517",
			"vol": "335.3697",
			"timestamp": 1703916551294
		}
	}`

	// モックのHTTPクライアント設定
	mockTransport := &MockHTTPTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := mockCryptoPrice
			r := io.NopCloser(bytes.NewReader([]byte(responseBody)))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r,
			}, nil
		},
	}

	mockHTTPClient := &http.Client{Transport: mockTransport}
	mockMarketPriceRepo := repoMarketPrice.NewCryptoRepository(mockHTTPClient)
	// オプションを使用してGraphQLサーバーをセットアップ
	opts := &graphql.SetupOptions{
		MockHTTPClient: mockHTTPClient,
		// MarketCryptoRepoにモックリポジトリを指定
		MarketCryptoRepo: mockMarketPriceRepo,
	}
    router := graphql.SetupGraphQLServer(db, opts)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.Crypto{Code: "btc", UserId: 1, Quantity: 0.05, GetPrice: 5047113.0})

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `query {
		cryptos{ id code getPrice quantity currentPrice }
	  }`
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            Cryptos []struct {
                ID string `json:"id"`
				Code string `json:"code"`
				GetPrice float64 `json:"getPrice"`
				Quantity float64 `json:"quantity"`
                CurrentPrice float64 `json:"currentPrice"`
            } `json:"cryptos"`
        } `json:"data"`
    }

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
    if len(response.Data.Cryptos) > 0 {
        assert.Equal(t, "btc", response.Data.Cryptos[0].Code)
		assert.Equal(t, 5047113.0, response.Data.Cryptos[0].GetPrice)
		assert.Equal(t, 0.05, response.Data.Cryptos[0].Quantity)
		assert.Equal(t, 5956517.0, response.Data.Cryptos[0].CurrentPrice)
    } else {
        t.Fatalf("Expected non-empty MarketPrice array")
    }
	// DB初期化
	db.Unscoped().Where("1=1").Delete(&model.Crypto{})
}


func TestCreateCryptoE2E(t *testing.T) {
	db := test.SetupTestDB()

	// モックの HTTP レスポンスを設定
	mockCryptoPrice := `{
		"success": 1,
		"data": {
			"sell": "88.225",
			"buy": "88.224",
			"open": "89.720",
			"high": "90.370",
			"low": "87.291",
			"last": "88.242",
			"vol": "9433775.5445",
			"timestamp": 1703916974315
		}
	}`

	// モックのHTTPクライアント設定
	mockTransport := &MockHTTPTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			responseBody := mockCryptoPrice
	
			r := io.NopCloser(bytes.NewReader([]byte(responseBody)))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r,
			}, nil
		},
	}

	mockHTTPClient := &http.Client{Transport: mockTransport}
	mockMarketPriceRepo := repoMarketPrice.NewCryptoRepository(mockHTTPClient)
	// オプションを使用してGraphQLサーバーをセットアップ
	opts := &graphql.SetupOptions{
		MockHTTPClient: mockHTTPClient,
		// MarketCryptoRepoにモックリポジトリを指定
		MarketCryptoRepo: mockMarketPriceRepo,
	}
    router := graphql.SetupGraphQLServer(db, opts)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `mutation {
		createCrypto(input: {
		  code: "xrp",
		  getPrice: 88.0,
		  quantity: 1
		}) {
		  id
		  code
		  getPrice
		  quantity
		  currentPrice
		}
	  }
	  `
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)
	  t.Log(w.Body)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            CreateCrypto struct {
                ID string `json:"id"`
				Code string `json:"code"`
				GetPrice float64 `json:"getPrice"`
				Quantity float64 `json:"quantity"`
                CurrentPrice float64 `json:"currentPrice"`
            } `json:"createCrypto"`
        } `json:"data"`
    }

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
	assert.Equal(t, "xrp", response.Data.CreateCrypto.Code)
	assert.Equal(t, 88.0, response.Data.CreateCrypto.GetPrice)
	assert.Equal(t, 1.0, response.Data.CreateCrypto.Quantity)
	assert.Equal(t, 88.242, response.Data.CreateCrypto.CurrentPrice)
	// DB初期化
	db.Unscoped().Where("1=1").Delete(&model.Crypto{})
}

func TestDeleteCryptoE2E(t *testing.T) {
	db := test.SetupTestDB()

	// テスト用データの追加
	crypto := model.Crypto{Code: "eth", UserId: 1, Quantity: 0.1, GetPrice: 200000.0}
	db.Create(&crypto)

	// 作成されたレコードのIDを取得
	createdCryptoID := crypto.ID

	// GraphQLサーバーのセットアップ
	router := graphql.SetupGraphQLServer(db, nil)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}
	// createdCryptoIDを文字列に変換
	createdCryptoIDStr := strconv.FormatUint(uint64(createdCryptoID), 10)

	// GraphQLリクエストの実行
	query := fmt.Sprintf(`mutation {
		deleteCrypto(id: "%s")
	}`, createdCryptoIDStr)

	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			DeleteCrypto bool `json:"deleteCrypto"`
		} `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.True(t, response.Data.DeleteCrypto)

	// データベースから削除されたことを確認
	var cryptoAfterDelete model.Crypto
	result := db.First(&cryptoAfterDelete, "id = ?", crypto.ID)
	assert.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
}
