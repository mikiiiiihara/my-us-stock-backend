package marketprice

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"my-us-stock-backend/repository/market-price/dto"
	"my-us-stock-backend/repository/market-price/entity"
	"net/http"
	"os"
	"strings"
	"time"
)

// MarketPriceRepository は仮想通貨の価格を取得するためのインターフェースです。
type MarketPriceRepository interface {
	FetchMarketPriceList(ctx context.Context, tickers []string) ([]dto.MarketPriceDto, error)
    FetchDividend(ctx context.Context, ticker string) (*entity.DividendEntity, error)
}

// DefaultMarketPriceRepository は CryptoRepository のデフォルト実装です。
type DefaultMarketPriceRepository struct {
    httpClient *http.Client
	baseURL string
	tickerToken string
	dividendMainToken string
	dividendSubToken string
}

// NewMarketPriceRepository は新しい DefaultCryptoRepository インスタンスを作成します。
func NewMarketPriceRepository(client *http.Client) *DefaultMarketPriceRepository {
    if client == nil {
        client = http.DefaultClient
    }
	baseURL := os.Getenv("MARKET_PRICE_URL")
	tickerToken := os.Getenv("MARKET_PRICE_TICKER_TOKEN")
	dividendMainToken := os.Getenv("MARKET_PRICE_DIVIDEND_MAIN_TOKEN")
	dividendSubToken := os.Getenv("MARKET_PRICE_DIVIDEND_SUB_TOKEN")
    return &DefaultMarketPriceRepository{
        httpClient: client,
        baseURL: baseURL,
		tickerToken: tickerToken,
		dividendMainToken: dividendMainToken,
		dividendSubToken: dividendSubToken,
    }
}

// FetchMarketPriceList fetches the current market prices for a list of tickers.
func (repo *DefaultMarketPriceRepository) FetchMarketPriceList(ctx context.Context, tickers []string) ([]dto.MarketPriceDto, error) {
    baseUrl := repo.baseURL
    token := repo.tickerToken
    url := fmt.Sprintf("%s/v3/quote-order/%s?apikey=%s", baseUrl, strings.Join(tickers, ","), token)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := repo.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("error fetching market prices: status %d", resp.StatusCode)
    }

    var prices []dto.MarketPriceResponse
    err = json.NewDecoder(resp.Body).Decode(&prices)
    if err != nil {
        return nil, err
    }

    // 市場にデータが存在するかチェック
    if len(prices) == 0 {
        return nil, fmt.Errorf("the specified tickers were not found")
    }

    // MarketPriceResponseからMarketPriceDtoに変換
    var priceDtos []dto.MarketPriceDto
    for _, price := range prices {
        priceDtos = append(priceDtos, dto.MarketPriceDto{
            Ticker:       price.Symbol,
            CurrentPrice: price.Price,
            PriceGets:    price.Change,
            CurrentRate:  price.ChangesPercentage,
        })
    }

    return priceDtos, nil
}


func (repo *DefaultMarketPriceRepository) FetchDividend(ctx context.Context, ticker string) (*entity.DividendEntity, error) {
    token := repo.dividendMainToken
    res, err := repo.fetchDividendApi(ctx, token, ticker)
    if err != nil {
        // 429 エラーの場合、別のトークンを使用して再試行
        if err.Error() == "rate limit exceeded" {
            token = repo.dividendSubToken
            res, err = repo.fetchDividendApi(ctx, token, ticker)
            if err != nil {
                return nil, fmt.Errorf("配当情報の取得に失敗しました。しばらく待ってからアクセスしてください: %w", err)
            }
        } else {
            return nil, fmt.Errorf("配当情報の取得に失敗しました: %w", err)
        }
    }
    return repo.createDividendEntity(res), nil
}

// FetchDividend は指定された銘柄の配当情報を取得します。
func (repo *DefaultMarketPriceRepository) fetchDividendApi(ctx context.Context, token string, ticker string) (*dto.DividendResponse, error) {
    url := fmt.Sprintf("%s/v3/historical-price-full/stock_dividend/%s?apikey=%s", repo.baseURL, ticker, token)
    resp, err := repo.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // ステータスコードのチェック
    if resp.StatusCode == 429 {
        return nil, fmt.Errorf("rate limit exceeded")
    }

    var dividendResponse dto.DividendResponse
    err = json.NewDecoder(resp.Body).Decode(&dividendResponse)
    if err != nil {
        return nil, err
    }

    return &dividendResponse, nil
}

// parseMonth は日付文字列から月を解析します。
func parseMonth(dateStr string) int {
    date, _ := time.Parse("2006-01-02", dateStr)
    return int(date.Month())
}

// roundToThreeDecimals は数値を小数点以下3桁に丸めます。
func roundToThreeDecimals(num float64) float64 {
	return math.Round(num*1000) / 1000
}