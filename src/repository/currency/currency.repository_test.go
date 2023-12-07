package currency

import (
	"bytes"
	"io" // ioutil の代わりに io をインポート
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
	mockResponseBody := `{"quotes":[{"currencyPairCode":"USDJPY","bid":"133.69"}]}`

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

	// テスト用の通貨URL
	testCurrencyURL := "http://test.url"

	// モックの HTTP クライアントとテスト用の URL を使用してリポジトリを初期化
	repo := NewCurrencyRepository(mockHTTPClient, testCurrencyURL)
	got, err := repo.FetchCurrentUsdJpy()

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, 133.69, got)
}
