package currency

import (
	"bytes"
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

func TestCurrencyRepository_FetchCurrentUsdJpy(t *testing.T) {
	// モックレスポンスの準備
	mockResponseBody := `{
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

	// モックのトランスポートを使用して HTTP クライアントを初期化
	mockTransport := &MockHTTPTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			// ioutil.NopCloser の代わりに io.NopCloser を使用
			r := io.NopCloser(bytes.NewReader([]byte(mockResponseBody)))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r,
			}, nil
		},
	}
	mockHTTPClient := &http.Client{Transport: mockTransport}

	// モックの HTTP クライアントとテスト用の URL を使用してリポジトリを初期化
	repo := NewCurrencyRepository(mockHTTPClient)
	got, err := repo.FetchCurrentUsdJpy()

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, 133.69, got)
}
