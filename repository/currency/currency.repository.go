package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"my-us-stock-backend/repository/currency/dto"
	"net/http"
	"os"
	"strconv"
)

type CurrencyRepository struct {
    httpClient *http.Client
    currencyURL string
}

func NewCurrencyRepository(client *http.Client) *CurrencyRepository {
    if client == nil {
        client = http.DefaultClient
    }
    currencyURL := os.Getenv("CURRENCY_URL")
    return &CurrencyRepository{
        httpClient: client,
        currencyURL: currencyURL,
    }
}

func (repo *CurrencyRepository) FetchCurrentUsdJpy() (float64, error) {
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