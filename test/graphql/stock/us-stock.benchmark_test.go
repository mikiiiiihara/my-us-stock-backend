package stock

import (
	// 必要なパッケージのインポート
	"bytes"
	"io"
	"my-us-stock-backend/app/database/model"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkUsStocks(b *testing.B) {
    db := test.SetupTestDB()

	// モックの HTTP レスポンスを設定
	mockStockPrice := `[
		{
			"symbol": "AAPL",
			"name": "Apple Inc.",
			"price": 193.94,
			"changesPercentage": 0.409,
			"change": 0.79,
			"dayLow": 193.59,
			"dayHigh": 194.6599,
			"yearHigh": 199.62,
			"yearLow": 124.17,
			"marketCap": 3016310032000,
			"priceAvg50": 185.9714,
			"priceAvg200": 179.08525,
			"exchange": "NASDAQ",
			"volume": 21883820,
			"avgVolume": 53390033,
			"open": 194.14,
			"previousClose": 193.15,
			"eps": 6.14,
			"pe": 31.59,
			"earningsAnnouncement": "2024-01-31T10:59:00.000+0000",
			"sharesOutstanding": 15552800000,
			"timestamp": 1703793087
		}
	]`

	mockDividend := `{
		"symbol": "AAPL",
		"historical": [
			{
				"date": "2023-11-10",
				"label": "November 10, 23",
				"adjDividend": 0.24,
				"dividend": 0.24,
				"recordDate": "2023-11-13",
				"paymentDate": "2023-11-16",
				"declarationDate": "2023-11-02"
			  },
			  {
				"date": "2023-08-11",
				"label": "August 11, 23",
				"adjDividend": 0.24,
				"dividend": 0.24,
				"recordDate": "2023-08-14",
				"paymentDate": "2023-08-17",
				"declarationDate": "2023-08-03"
			  },
			  {
				"date": "2023-05-12",
				"label": "May 12, 23",
				"adjDividend": 0.24,
				"dividend": 0.24,
				"recordDate": "2023-05-15",
				"paymentDate": "2023-05-18",
				"declarationDate": "2023-05-04"
			  },
			  {
				"date": "2023-02-10",
				"label": "February 10, 23",
				"adjDividend": 0.23,
				"dividend": 0.23,
				"recordDate": "2023-02-13",
				"paymentDate": "2023-02-16",
				"declarationDate": "2023-02-02"
			  },
			  {
				"date": "2022-11-04",
				"label": "November 04, 22",
				"adjDividend": 0.23,
				"dividend": 0.23,
				"recordDate": "2022-11-07",
				"paymentDate": "2022-11-10",
				"declarationDate": "2022-10-27"
			  }
		]
	}
	`

	// モックのHTTPクライアント設定
	mockTransport := &MockHTTPTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			var responseBody string
	
			// URLに基づいて異なるレスポンスを返す
			if req.URL.Path == "/v3/quote-order/AAPL" {
				responseBody = mockStockPrice
			} else if req.URL.Path == "/v3/historical-price-full/stock_dividend/AAPL" {
				responseBody = mockDividend
			}
	
			r := io.NopCloser(bytes.NewReader([]byte(responseBody)))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       r,
			}, nil
		},
	}

	mockHTTPClient := &http.Client{Transport: mockTransport}
	mockMarketPriceRepo := repoMarketPrice.NewMarketPriceRepository(mockHTTPClient)
	// オプションを使用してGraphQLサーバーをセットアップ
	opts := &graphql.SetupOptions{
		MockHTTPClient: mockHTTPClient,
		// MarketPriceRepoにモックリポジトリを指定
		MarketPriceRepo: mockMarketPriceRepo,
	}

    router := graphql.SetupGraphQLServer(db, opts)
    ts := httptest.NewServer(router)
    defer ts.Close()

	// テスト用のユーザーを作成
	db.Create(&model.UsStock{Code: "AAPL", UserId: 1, Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9})


    token, err := graphql.GenerateTestAccessTokenForUserId(1)
    if err != nil {
        b.Fatalf("Failed to generate test access token: %v", err)
    }

    query := `query {
        usStocks{ id code getPrice sector dividend quantity usdJpy currentPrice priceGets currentRate }
    }`

    // ベンチマークの実行
    b.ResetTimer() // タイマーをリセットしてから計測開始
    for i := 0; i < b.N; i++ {
        graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)
    }
}
