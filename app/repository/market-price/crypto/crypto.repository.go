package crypto

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// CryptoRepository は仮想通貨の価格を取得するためのインターフェースです。
type CryptoRepository interface {
    FetchCryptoPrice(code string) (*Crypto, error)
}

// DefaultCryptoRepository は CryptoRepository のデフォルト実装です。
type DefaultCryptoRepository struct {
    httpClient *http.Client
    cryptoURL string
}

// NewCryptoRepository は新しい DefaultCryptoRepository インスタンスを作成します。
func NewCryptoRepository(client *http.Client) *DefaultCryptoRepository {
    if client == nil {
        client = http.DefaultClient
    }
	cryptoURL := os.Getenv("CRYPTO_URL")
    return &DefaultCryptoRepository{
        httpClient: client,
        cryptoURL: cryptoURL,
    }
}

// FetchCryptoPrice は指定された仮想通貨の価格を取得します。
func (repo *DefaultCryptoRepository) FetchCryptoPrice(code string) (*Crypto, error) {
    url := fmt.Sprintf("%s/%s_jpy/ticker", repo.cryptoURL, code)
    resp, err := repo.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var apiResp ApiResponse
    err = json.Unmarshal(body, &apiResp)
    if err != nil {
        return nil, err
    }

    price, err := strconv.ParseFloat(apiResp.Data.Last, 64)
    if err != nil {
        return nil, err
    }

    return &Crypto{
        Name:  string(code),
        Price: price,
    }, nil
}
