package marketprice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockHTTPTransport は http.RoundTripper のインターフェースを満たすモック実装です。
type MockHTTPTransport struct {
    RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip は http.RoundTripper の RoundTrip メソッドを模倣します。
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.RoundTripFunc(req)
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

func TestMarketPriceE2E(t *testing.T) {
    // モックの HTTP レスポンスを設定
    mockResponseBody := `[
        {
            "symbol": "AAPL",
            "price": 189.84,
            "changesPercentage": 0.0685,
            "change": 0.13
        },
        {
            "symbol": "KO",
            "price": 57.205,
            "changesPercentage": 0.0962,
            "change": 0.055
        }
    ]`

    // モックのHTTPクライアント設定
    mockTransport := &MockHTTPTransport{
        RoundTripFunc: func(req *http.Request) (*http.Response, error) {
            r := io.NopCloser(bytes.NewReader([]byte(mockResponseBody)))
            return &http.Response{
                StatusCode: http.StatusOK,
                Body:       r,
            }, nil
        },
    }
    mockHTTPClient := &http.Client{Transport: mockTransport}
    mockMarketPriceRepo := repoMarketPrice.NewMarketPriceRepository(mockHTTPClient)
    // オプションを使用してGraphQLサーバーをセットアップ
    opts := &graphql.SetupOptions{
        MockHTTPClient: mockHTTPClient,
        // MarketPriceRepoにモックリポジトリを指定
        MarketPriceRepo: mockMarketPriceRepo,
    }
    graphqlServer := graphql.SetupGraphQLServer(nil,opts)

	// GraphQLリクエストの実行
	query := `query {
		marketPrices(tickerList: ["AAPL","KO"]){
		  ticker
		  currentPrice
		  currentRate
		  priceGets
		}
	  }`
	w := executeGraphQLRequest(graphqlServer, query)

    // レスポンスのログ出力（デバッグ用）
    fmt.Printf("Response Body: %s\n", w.Body.String())

	// レスポンスボディの解析
    var response struct {
        Data struct {
            MarketPrices []struct {
                Ticker string `json:"ticker"`
                CurrentPrice float64 `json:"currentPrice"`
                PriceGets float64 `json:"priceGets"`
                CurrentRate float64 `json:"currentRate"`
            } `json:"marketPrices"`
        } `json:"data"`
    }

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
    if len(response.Data.MarketPrices) > 0 {
        assert.Equal(t, "AAPL", response.Data.MarketPrices[0].Ticker)
    } else {
        t.Fatalf("Expected non-empty MarketPrice array")
    }
}