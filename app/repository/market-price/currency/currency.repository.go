package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"my-us-stock-backend/app/repository/market-price/currency/dto"
	"net/http"
	"os"
	"strconv"
)

type CurrencyRepository interface {
    FetchCurrentUsdJpy(ctx context.Context) (float64, error)
}

// DefaultCurrencyRepository は CurrencyRepository のデフォルトの実装です。
type DefaultCurrencyRepository struct {
    httpClient *http.Client
    currencyURL string
}

// NewCurrencyRepository は新しい DefaultCurrencyRepository のインスタンスを作成します。
func NewCurrencyRepository(client *http.Client) *DefaultCurrencyRepository {
    if client == nil {
        client = http.DefaultClient
    }
    currencyURL := os.Getenv("CURRENCY_URL")
    return &DefaultCurrencyRepository{
        httpClient: client,
        currencyURL: currencyURL,
    }
}

func (repo *DefaultCurrencyRepository) FetchCurrentUsdJpy(ctx context.Context) (float64, error) {
    resp, err := repo.httpClient.Get(repo.currencyURL)
    if err != nil {
        log.Printf("Error fetching currency data: %v\n", err)
        return 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v\n", err)
        return 0, err
    }

    var fx dto.Fx
    err = json.Unmarshal(body, &fx)
    if err != nil {
        log.Printf("Error unmarshalling JSON: %v\n", err)
        return 0, err
    }

    for _, quote := range fx.Quotes {
        if quote.CurrencyPairCode == "USDJPY" {
            currentUsdJpy, err := strconv.ParseFloat(quote.Bid, 64)
            if err != nil {
                log.Printf("Error parsing float value: %v\n", err)
                return 0, err
            }
            return currentUsdJpy, nil
        }
    }

    log.Println("USDJPY not found in quotes")
    return 0, fmt.Errorf("USDJPY not found")
}
