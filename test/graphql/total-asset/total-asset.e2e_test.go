package totalasset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-us-stock-backend/app/database/model"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	repoMarketCrypto "my-us-stock-backend/app/repository/market-price/crypto"
	"my-us-stock-backend/app/repository/market-price/currency"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTotalAssetsE2E(t *testing.T) {
	db := test.SetupTestDB()
	router := graphql.SetupGraphQLServer(db, nil)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// テスト用データの追加
	db.Create(&model.TotalAsset{UserId: 1, CashJpy: 10000, CashUsd: 100, Stock: 50000})

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `query {
		totalAssets(day: 30) {
			id cashJpy cashUsd stock fund crypto fixedIncomeAsset createdAt
		}
	}`
	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			TotalAssets []struct {
				ID               string  `json:"id"`
				CashJpy          float64 `json:"cashJpy"`
				CashUsd          float64 `json:"cashUsd"`
				Stock            float64 `json:"stock"`
				Fund             float64 `json:"fund"`
				Crypto           float64 `json:"crypto"`
				FixedIncomeAsset float64 `json:"fixedIncomeAsset"`
				CreatedAt        string  `json:"createdAt"`
			} `json:"totalAssets"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	// レスポンスボディの内容の検証
	assert.NotEmpty(t, response.Data.TotalAssets)
	for _, asset := range response.Data.TotalAssets {
		assert.Greater(t, asset.CashJpy, 0.0)
		assert.Greater(t, asset.CashUsd, 0.0)
		assert.Greater(t, asset.Stock, 0.0)
	}
}


// MockHTTPTransport は http.RoundTripper のインターフェースを満たすモック実装です。
type MockHTTPTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip は http.RoundTripper の RoundTrip メソッドを模倣します。
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
func TestUpdateTotalAssetE2E(t *testing.T) {
	db := test.SetupTestDB()

	// モックの HTTP レスポンスを設定
	mockStockPrice := `[
		{
			"symbol": "AAPL",
			"name": "Apple Inc.",
			"price": 193.94,
			"changesPercentage": 0.409,
			"change": 0.79,
			"dayLow": 193.59,
			"dayHigh": 194.6599,
			"yearHigh": 199.62,
			"yearLow": 124.17,
			"marketCap": 3016310032000,
			"priceAvg50": 185.9714,
			"priceAvg200": 179.08525,
			"exchange": "NASDAQ",
			"volume": 21883820,
			"avgVolume": 53390033,
			"open": 194.14,
			"previousClose": 193.15,
			"eps": 6.14,
			"pe": 31.59,
			"earningsAnnouncement": "2024-01-31T10:59:00.000+0000",
			"sharesOutstanding": 15552800000,
			"timestamp": 1703793087
		}
	]`

	mockCurrency := `{
		"quotes": [
		  {
			"high": "1.2108",
			"open": "1.2093",
			"bid": "1.2105",
			"currencyPairCode": "GBPUSD",
			"ask": "1.2115",
			"low": "1.2091"
		  },
		  {
			"high": "133.74",
			"open": "133.73",
			"bid": "133.69",
			"currencyPairCode": "USDJPY",
			"ask": "133.72",
			"low": "133.69"
		  },
		  {
			"high": "1.5938",
			"open": "1.5936",
			"bid": "1.5936",
			"currencyPairCode": "EURAUD",
			"ask": "1.5949",
			"low": "1.5923"
		  }
		]
	  }`

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
			var responseBody string
			// URLに基づいて異なるレスポンスを返す
			if req.URL.Path == "/v3/quote-order/AAPL" {
				responseBody = mockStockPrice
			} else if req.URL.Path == "" {
				responseBody = mockCurrency
			} else if req.URL.Path == "/btc_jpy/ticker" {
				responseBody = mockCryptoPrice
			}
	
			r := io.NopCloser(bytes.NewReader([]byte(responseBody)))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r,
			}, nil
		},
	}
	mockHTTPClient := &http.Client{Transport: mockTransport}
	mockMarketPriceRepo := repoMarketPrice.NewMarketPriceRepository(mockHTTPClient)
	mockCurrencyRepo := currency.NewCurrencyRepository(mockHTTPClient)
	mockMarketCryptoRepo := repoMarketCrypto.NewCryptoRepository(mockHTTPClient)
	// オプションを使用してGraphQLサーバーをセットアップ
	opts := &graphql.SetupOptions{
		MockHTTPClient: mockHTTPClient,
		// MarketPriceRepoにモックリポジトリを指定
		MarketPriceRepo: mockMarketPriceRepo,
		CurrencyRepo: mockCurrencyRepo,
		MarketCryptoRepo: mockMarketCryptoRepo,
	}
    router := graphql.SetupGraphQLServer(db, opts)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// テスト用データの追加
	totalAsset := model.TotalAsset{UserId: 1, CashJpy: 10000, CashUsd: 100, Stock: 50000}
	db.Create(&totalAsset)
	db.Create(&model.UsStock{Code: "AAPL", UserId: 1, Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9})
	db.Create(&model.FixedIncomeAsset{Code: "Funds", UserId: 1, DividendRate: 3.5, GetPriceTotal: 100000.0, PaymentMonth: pq.Int64Array{6, 12}})
    db.Create(&model.JapanFund{Code: "SP500", Name:"ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", GetPrice: 15523.81, GetPriceTotal: 761157.0,UserId: 1})
    db.Create(&model.Crypto{Code: "btc", UserId: 1, Quantity: 0.05, GetPrice: 5047113.0})

	// 作成されたレコードのIDを取得
	createdTotalAssetID := totalAsset.ID

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// createdTotalAssetIDを文字列に変換
	createdTotalAssetIDStr := strconv.FormatUint(uint64(createdTotalAssetID), 10)

	// GraphQLリクエストの実行
	query := fmt.Sprintf(`mutation {
		updateTotalAsset(input: {
			id: "%s"
			cashJpy: 15000
			cashUsd: 150
		}) {
			id cashJpy cashUsd stock fund crypto fixedIncomeAsset createdAt
		}
	}`, createdTotalAssetIDStr)
	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			UpdateTotalAsset struct {
				ID               string  `json:"id"`
				CashJpy          float64 `json:"cashJpy"`
				CashUsd          float64 `json:"cashUsd"`
				Stock            float64 `json:"stock"`
				Fund             float64 `json:"fund"`
				Crypto           float64 `json:"crypto"`
				FixedIncomeAsset float64 `json:"fixedIncomeAsset"`
				CreatedAt        string  `json:"createdAt"`
			} `json:"updateTotalAsset"`
		} `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	// レスポンスボディの内容の検証
	assert.Equal(t, createdTotalAssetIDStr, response.Data.UpdateTotalAsset.ID)
	assert.Equal(t, 15000.0, response.Data.UpdateTotalAsset.CashJpy)
	assert.Equal(t, 150.0, response.Data.UpdateTotalAsset.CashUsd)
	assert.Equal(t, 259278.0, response.Data.UpdateTotalAsset.Stock)
	assert.Equal(t, 1193527.0, response.Data.UpdateTotalAsset.Fund)
	assert.Equal(t, 297826.0, response.Data.UpdateTotalAsset.Crypto)
	assert.Equal(t, 100000.0, response.Data.UpdateTotalAsset.FixedIncomeAsset)
}
