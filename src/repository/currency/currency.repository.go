package currency

import (
	"encoding/json"
	"fmt"
	"io"
	"my-us-stock-backend/src/repository/currency/dto"
	"net/http"
	"os"
	"strconv"
)

type CurrencyRepository struct {
	httpClient *http.Client
}

func NewCurrencyRepository(client *http.Client) *CurrencyRepository {
	return &CurrencyRepository{
		httpClient: client,
	}
}

func (repo *CurrencyRepository) FetchCurrentUsdJpy() (float64, error) {
	currencyURL := os.Getenv("CURRENCY_URL")
	resp, err := repo.httpClient.Get(currencyURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var fx dto.Fx
	err = json.Unmarshal(body, &fx)
	if err != nil {
		return 0, err
	}

	for _, quote := range fx.Quotes {
		if quote.CurrencyPairCode == "USDJPY" {
			currentUsdJpy, err := strconv.ParseFloat(quote.Bid, 64)
			if err != nil {
				return 0, err
			}
			return currentUsdJpy, nil
		}
	}

	return 0, fmt.Errorf("USDJPY not found")
}
