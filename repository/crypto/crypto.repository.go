package crypto

import (
	"encoding/json"
	"fmt"
	"io"
	"my-us-stock-backend/repository/crypto/dto"
	"net/http"
	"os"
	"strconv"
)

type CryptoRepository struct {
    httpClient *http.Client
    cryptoURL string
}

func NewCryptoRepository(client *http.Client) *CryptoRepository {
    if client == nil {
        client = http.DefaultClient
    }
    cryptoURL := os.Getenv("CRYPTO_URL")
    return &CryptoRepository{
        httpClient: client,
        cryptoURL: cryptoURL,
    }
}

// FetchCryptoPrice fetches the price of the given cryptocurrency
func (repo *CryptoRepository) FetchCryptoPrice(ticker CryptoTicker) (*dto.Crypto, error) {
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