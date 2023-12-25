package crypto

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

func TestCryptoRepository_FetchCryptoPrice(t *testing.T) {
	// モックレスポンスの準備
	mockResponseBody := `{
		"success": 1,
		"data": {
			"sell": "50.750",
			"buy": "50.749",
			"open": "50.706",
			"high": "50.917",
			"low": "49.333",
			"last": "50.749",
			"vol": "13346627.3932",
			"timestamp": 1679376127932
		}
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

	// リポジトリのインスタンスを作成
	repo := NewCryptoRepository(mockHTTPClient)

	// テストの実行
	t.Run("正常に仮想通貨の価格を取得", func(t *testing.T) {
		expected := &Crypto{
			Name:  "btc",
			Price: 50.749,
		}

		result, err := repo.FetchCryptoPrice("btc")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
