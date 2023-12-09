package marketprice

import (
	"bytes"
	"context"
	"io"
	"net/http"
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

func TestFetchMarketPriceList(t *testing.T) {
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

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	tickers := []string{"AAPL", "KO"}
	prices, err := repo.FetchMarketPriceList(context.Background(), tickers)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, prices, 2)
	// AAPL
	assert.Equal(t, "AAPL", prices[0].Ticker)
	assert.Equal(t, 189.84, prices[0].CurrentPrice)
	assert.Equal(t, 0.0685, prices[0].CurrentRate)
	assert.Equal(t, 0.13, prices[0].PriceGets)
	// KO
	assert.Equal(t, "KO", prices[1].Ticker)
	assert.Equal(t, 57.205, prices[1].CurrentPrice)
	assert.Equal(t, 0.0685, prices[0].CurrentRate)
	assert.Equal(t, 0.13, prices[0].PriceGets)
}

func TestFetchMarketPriceList_NoData(t *testing.T) {
	// モックの HTTP レスポンスを設定（空のリスト）
	mockResponseBody := `[]`

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	tickers := []string{"AAPLKK"}
	_, err := repo.FetchMarketPriceList(context.Background(), tickers)

	// アサーション：期待されるエラーメッセージをチェック
	expectedErrorMessage := "the specified tickers were not found" // 期待されるエラーメッセージに更新
	assert.Equal(t, expectedErrorMessage, err.Error())
}