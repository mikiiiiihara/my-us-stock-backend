package marketprice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	serviceCurrency "my-us-stock-backend/graphql/currency"
	"my-us-stock-backend/graphql/generated"
	serviceMarketPrice "my-us-stock-backend/graphql/market-price"
	serviceUser "my-us-stock-backend/graphql/user"
	repoCurrency "my-us-stock-backend/repository/currency"
	repoMarketPrice "my-us-stock-backend/repository/market-price"
	repoUser "my-us-stock-backend/repository/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
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

type CustomQueryResolver struct {
	UserResolver     *serviceUser.Resolver
	CurrencyResolver *serviceCurrency.Resolver
	MarketPriceResolver *serviceMarketPrice.Resolver
}


func (r *CustomQueryResolver) Query() generated.QueryResolver {
    return r
}

func (r *CustomQueryResolver) Mutation() generated.MutationResolver {
    // ここでuserResolverを使用してMutationを実装する
    return r.UserResolver
}

func (r *CustomQueryResolver) User(ctx context.Context, id string) (*generated.User, error) {
    return r.UserResolver.User(ctx, id)
}

func (r *CustomQueryResolver) GetCurrentUsdJpy(ctx context.Context) (float64, error) {
    return r.CurrencyResolver.GetCurrentUsdJpy(ctx)
}

// GetMarketPricesメソッドの実装
func (r *CustomQueryResolver) GetMarketPrices(ctx context.Context, tickers []*string) ([]*generated.MarketPrice, error) {
    // 文字列スライスに変換
    tickerStrs := make([]string, len(tickers))
    for i, t := range tickers {
        tickerStrs[i] = *t
    }

    // サービスを呼び出して結果を取得
    return r.MarketPriceResolver.GetMarketPrices(ctx, tickerStrs)
}

func NewCustomQueryResolver(userResolver *serviceUser.Resolver, currencyResolver *serviceCurrency.Resolver, marketPriceResolver *serviceMarketPrice.Resolver) *CustomQueryResolver {
    return &CustomQueryResolver{
        UserResolver:     userResolver,
        CurrencyResolver: currencyResolver,
		MarketPriceResolver: marketPriceResolver,
    }
}

// setupGraphQLServer はテスト用のGraphQLサーバーをセットアップします
func setupGraphQLServer(mockHTTPClient *http.Client) *handler.Server {
    // リポジトリ、サービス、リゾルバの初期化
	currencyRepo := repoCurrency.NewCurrencyRepository(nil)
    currencyService := serviceCurrency.NewCurrencyService(currencyRepo)
	currencyResolver := serviceCurrency.NewResolver(currencyService)

    userRepo := repoUser.NewUserRepository(nil)
    userService := serviceUser.NewUserService(userRepo)
    userResolver := serviceUser.NewResolver(userService)

	marketPriceRepo := repoMarketPrice.NewMarketPriceRepository(mockHTTPClient)
	marketPriceService := serviceMarketPrice.NewMarketPriceService(marketPriceRepo)
    marketPriceResolver := serviceMarketPrice.NewResolver(marketPriceService)

// CustomQueryResolverの初期化
resolver := NewCustomQueryResolver(userResolver, currencyResolver, marketPriceResolver)

return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))
}


// executeGraphQLRequest はGraphQLリクエストを実行し、レスポンスを返します
func executeGraphQLRequest(h http.Handler, query string) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
	})
	req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestMarketPriceE2E(t *testing.T) {
    // モックの HTTP レスポンスを設定
    mockResponseBody := `[
        {
            "symbol": "AAPL",
            "price": 189.84,
            "changesPercentage": 0.0685,
            "change": 0.13
        },
        {
            "symbol": "KO",
            "price": 57.205,
            "changesPercentage": 0.0962,
            "change": 0.055
        }
    ]`

    // モックのHTTPクライアント設定
    mockTransport := &MockHTTPTransport{
        RoundTripFunc: func(req *http.Request) (*http.Response, error) {
            r := io.NopCloser(bytes.NewReader([]byte(mockResponseBody)))
            return &http.Response{
                StatusCode: http.StatusOK,
                Body:       r,
            }, nil
        },
    }
    mockHTTPClient := &http.Client{Transport: mockTransport}

    // GraphQLサーバーのセットアップ
    graphqlServer := setupGraphQLServer(mockHTTPClient)

	// GraphQLリクエストの実行
	query := `query {
		getMarketPrices(tickerList: ["AAPL","KO"]){
		  ticker
		  currentPrice
		  currentRate
		  priceGets
		}
	  }`
	w := executeGraphQLRequest(graphqlServer, query)

    // レスポンスのログ出力（デバッグ用）
    fmt.Printf("Response Body: %s\n", w.Body.String())

	// レスポンスボディの解析
    var response struct {
        Data struct {
            MarketPrices []struct {
                Ticker string `json:"ticker"`
                CurrentPrice float64 `json:"currentPrice"`
                PriceGets float64 `json:"priceGets"`
                CurrentRate float64 `json:"currentRate"`
            } `json:"getMarketPrices"`
        } `json:"data"`
    }

	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
    if len(response.Data.MarketPrices) > 0 {
        assert.Equal(t, "AAPL", response.Data.MarketPrices[0].Ticker)
    } else {
        t.Fatalf("Expected non-empty MarketPrice array")
    }
}