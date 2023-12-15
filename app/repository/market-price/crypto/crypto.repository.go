package crypto

import (
	"encoding/json"
	"fmt"
	"io"
	"my-us-stock-backend/app/repository/market-price/crypto/dto"
	"net/http"
	"os"
	"strconv"
)

// CryptoRepository は仮想通貨の価格を取得するためのインターフェースです。
type CryptoRepository interface {
    FetchCryptoPrice(ticker CryptoTicker) (*dto.Crypto, error)
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
func (repo *DefaultCryptoRepository) FetchCryptoPrice(ticker CryptoTicker) (*dto.Crypto, error) {
    url := fmt.Sprintf("%s/%s_jpy/ticker", repo.cryptoURL, ticker)
    resp, err := repo.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var apiResp dto.ApiResponse
    err = json.Unmarshal(body, &apiResp)
    if err != nil {
        return nil, err
    }

    price, err := strconv.ParseFloat(apiResp.Data.Last, 64)
    if err != nil {
        return nil, err
    }

    return &dto.Crypto{
        Name:  string(ticker),
        Price: price,
    }, nil
}
