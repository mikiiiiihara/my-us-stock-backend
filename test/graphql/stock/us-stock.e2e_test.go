package stock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"my-us-stock-backend/app/database/model"
	repoMarketPrice "my-us-stock-backend/app/repository/market-price"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// MockHTTPTransport は http.RoundTripper のインターフェースを満たすモック実装です。
type MockHTTPTransport struct {
    RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip は http.RoundTripper の RoundTrip メソッドを模倣します。
func (m *MockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.RoundTripFunc(req)
}

func TestUsStocksE2E(t *testing.T) {
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

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.UsStock{Code: "AAPL", UserId: 1, Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9})

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `query {
		usStocks{ id code getPrice sector dividend quantity usdJpy currentPrice priceGets currentRate }
	  }`
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            UsStocks []struct {
                ID string `json:"id"`
				Code string `json:"code"`
				GetPrice float64 `json:"getPrice"`
				Sector string `json:"sector"`
				Dividend float64 `json:"dividend"`
				Quantity float64 `json:"quantity"`
				UsdJpy float64 `json:"usdJpy"`
                CurrentPrice float64 `json:"currentPrice"`
                PriceGets float64 `json:"priceGets"`
                CurrentRate float64 `json:"currentRate"`
            } `json:"usStocks"`
        } `json:"data"`
    }

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
    if len(response.Data.UsStocks) > 0 {
        assert.Equal(t, "AAPL", response.Data.UsStocks[0].Code)
		assert.Equal(t, 100.0, response.Data.UsStocks[0].GetPrice)
		assert.Equal(t, "IT", response.Data.UsStocks[0].Sector)
		assert.Equal(t, 0.95, response.Data.UsStocks[0].Dividend)
		assert.Equal(t, 10.0, response.Data.UsStocks[0].Quantity)
		assert.Equal(t, 133.9, response.Data.UsStocks[0].UsdJpy)
		assert.Equal(t, 193.94, response.Data.UsStocks[0].CurrentPrice)
		assert.Equal(t, 0.79, response.Data.UsStocks[0].PriceGets)
		assert.Equal(t, 0.409, response.Data.UsStocks[0].CurrentRate)
    } else {
        t.Fatalf("Expected non-empty MarketPrice array")
    }
}


func TestCreateUsStockE2E(t *testing.T) {
	db := test.SetupTestDB()

	// モックの HTTP レスポンスを設定
	mockStockPrice := `[
		{
			"symbol": "VTI",
			"name": "Vanguard Total Stock Market Index Fund",
			"price": 238.55,
			"changesPercentage": 0.1259,
			"change": 0.3,
			"dayLow": 238.15,
			"dayHigh": 238.7399,
			"yearHigh": 238.7399,
			"yearLow": 188.93,
			"marketCap": 494685928424,
			"priceAvg50": 222.1308,
			"priceAvg200": 215.4878,
			"exchange": "AMEX",
			"volume": 3124566,
			"avgVolume": 3536133,
			"open": 238.25,
			"previousClose": 238.25,
			"eps": 10.632924,
			"pe": 22.44,
			"earningsAnnouncement": null,
			"sharesOutstanding": 2073720094,
			"timestamp": 1703793061
		}
	]`

	mockDividend := `{
		"symbol": "VTI",
		"historical": [
			{
				"date": "2023-12-21",
				"label": "December 21, 23",
				"adjDividend": 1.002,
				"dividend": 1.002,
				"recordDate": "",
				"paymentDate": "",
				"declarationDate": ""
			},
			{
				"date": "2023-09-21",
				"label": "September 21, 23",
				"adjDividend": 0.7984,
				"dividend": 0.7984,
				"recordDate": "2023-09-22",
				"paymentDate": "2023-09-26",
				"declarationDate": "2023-03-17"
			},
			{
				"date": "2023-06-23",
				"label": "June 23, 23",
				"adjDividend": 0.8265,
				"dividend": 0.8265,
				"recordDate": "2023-06-26",
				"paymentDate": "2023-06-28",
				"declarationDate": "2023-03-17"
			},
			{
				"date": "2023-03-23",
				"label": "March 23, 23",
				"adjDividend": 0.7862,
				"dividend": 0.7862,
				"recordDate": "2023-03-24",
				"paymentDate": "2023-03-28",
				"declarationDate": "2023-03-17"
			},
			{
				"date": "2022-12-22",
				"label": "December 22, 22",
				"adjDividend": 0.9305,
				"dividend": 0.931,
				"recordDate": "2022-12-23",
				"paymentDate": "2022-12-28",
				"declarationDate": "2022-12-20"
			},
			{
				"date": "2022-09-23",
				"label": "September 23, 22",
				"adjDividend": 0.7955,
				"dividend": 0.796,
				"recordDate": "2022-09-26",
				"paymentDate": "2022-09-28",
				"declarationDate": "2022-09-21"
			}
		]
	}
	`

	// モックのHTTPクライアント設定
	mockTransport := &MockHTTPTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			var responseBody string
	
			// URLに基づいて異なるレスポンスを返す
			if req.URL.Path == "/v3/quote-order/VTI" {
				responseBody = mockStockPrice
			} else if req.URL.Path == "/v3/historical-price-full/stock_dividend/VTI" {
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

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `mutation {
		createUsStock(input: {
		  code: "VTI",
		  getPrice: 190.0,
		  quantity: 10.0,
		  sector: "Index",
		  usdJpy: 130.0
		}) {
		  id
		  code
		  getPrice
		  dividend
		  quantity
		  sector
		  usdJpy
		  currentPrice
		  priceGets
		  currentRate
		}
	  }
	  `
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)
	  t.Log(w.Body)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            CreateUsStock struct {
                ID string `json:"id"`
				Code string `json:"code"`
				GetPrice float64 `json:"getPrice"`
				Sector string `json:"sector"`
				Dividend float64 `json:"dividend"`
				Quantity float64 `json:"quantity"`
				UsdJpy float64 `json:"usdJpy"`
                CurrentPrice float64 `json:"currentPrice"`
                PriceGets float64 `json:"priceGets"`
                CurrentRate float64 `json:"currentRate"`
            } `json:"createUsStock"`
        } `json:"data"`
    }

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
	assert.Equal(t, "VTI", response.Data.CreateUsStock.Code)
	assert.Equal(t, 190.0, response.Data.CreateUsStock.GetPrice)
	assert.Equal(t, "Index", response.Data.CreateUsStock.Sector)
	assert.Equal(t, 2.411, response.Data.CreateUsStock.Dividend)
	assert.Equal(t, 10.0, response.Data.CreateUsStock.Quantity)
	assert.Equal(t, 130.0, response.Data.CreateUsStock.UsdJpy)
	assert.Equal(t, 238.55, response.Data.CreateUsStock.CurrentPrice)
	assert.Equal(t, 0.3, response.Data.CreateUsStock.PriceGets)
	assert.Equal(t, 0.1259, response.Data.CreateUsStock.CurrentRate)
}

func TestUpdateUsStockE2E(t *testing.T) {
	db := test.SetupTestDB()

	// モックの HTTP レスポンスを設定
	mockStockPrice := `[
		{
			"symbol": "VTI",
			"name": "Vanguard Total Stock Market Index Fund",
			"price": 238.55,
			"changesPercentage": 0.1259,
			"change": 0.3,
			"dayLow": 238.15,
			"dayHigh": 238.7399,
			"yearHigh": 238.7399,
			"yearLow": 188.93,
			"marketCap": 494685928424,
			"priceAvg50": 222.1308,
			"priceAvg200": 215.4878,
			"exchange": "AMEX",
			"volume": 3124566,
			"avgVolume": 3536133,
			"open": 238.25,
			"previousClose": 238.25,
			"eps": 10.632924,
			"pe": 22.44,
			"earningsAnnouncement": null,
			"sharesOutstanding": 2073720094,
			"timestamp": 1703793061
		}
	]`

	mockDividend := `{
		"symbol": "VTI",
		"historical": [
			{
				"date": "2023-12-21",
				"label": "December 21, 23",
				"adjDividend": 1.002,
				"dividend": 1.002,
				"recordDate": "",
				"paymentDate": "",
				"declarationDate": ""
			},
			{
				"date": "2023-09-21",
				"label": "September 21, 23",
				"adjDividend": 0.7984,
				"dividend": 0.7984,
				"recordDate": "2023-09-22",
				"paymentDate": "2023-09-26",
				"declarationDate": "2023-03-17"
			},
			{
				"date": "2023-06-23",
				"label": "June 23, 23",
				"adjDividend": 0.8265,
				"dividend": 0.8265,
				"recordDate": "2023-06-26",
				"paymentDate": "2023-06-28",
				"declarationDate": "2023-03-17"
			},
			{
				"date": "2023-03-23",
				"label": "March 23, 23",
				"adjDividend": 0.7862,
				"dividend": 0.7862,
				"recordDate": "2023-03-24",
				"paymentDate": "2023-03-28",
				"declarationDate": "2023-03-17"
			},
			{
				"date": "2022-12-22",
				"label": "December 22, 22",
				"adjDividend": 0.9305,
				"dividend": 0.931,
				"recordDate": "2022-12-23",
				"paymentDate": "2022-12-28",
				"declarationDate": "2022-12-20"
			},
			{
				"date": "2022-09-23",
				"label": "September 23, 22",
				"adjDividend": 0.7955,
				"dividend": 0.796,
				"recordDate": "2022-09-26",
				"paymentDate": "2022-09-28",
				"declarationDate": "2022-09-21"
			}
		]
	}
	`

	// モックのHTTPクライアント設定
	mockTransport := &MockHTTPTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			var responseBody string
	
			// URLに基づいて異なるレスポンスを返す
			if req.URL.Path == "/v3/quote-order/VTI" {
				responseBody = mockStockPrice
			} else if req.URL.Path == "/v3/historical-price-full/stock_dividend/VTI" {
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

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用データの追加
    usStock := model.UsStock{Code: "VTI", UserId: 1, Quantity: 10, GetPrice: 230, Sector: "Index", UsdJpy: 133.9}
    db.Create(&usStock)

    // 更新前の株式情報を取得
    createdUsStockID := usStock.ID

    // ダミーのアクセストークンを生成
    token, err := graphql.GenerateTestAccessTokenForUserId(1)
    if err != nil {
        t.Fatalf("Failed to generate test access token: %v", err)
    }

    // GraphQLリクエストの実行
    updateQuery := fmt.Sprintf(`mutation {
        updateUsStock(input: {
            id: "%s",
            getPrice: 234.0,
            quantity: 15,
            usdJpy: 135.0
        }) {
            id
            code
            getPrice
            quantity
            sector
            usdJpy
            currentPrice
            priceGets
            currentRate
        }
    }`, strconv.FormatUint(uint64(createdUsStockID), 10))

    w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, updateQuery, token)

    // レスポンスボディの解析
    var response struct {
        Data struct {
            UpdateUsStock struct {
                ID          string  `json:"id"`
                Code        string  `json:"code"`
                GetPrice    float64 `json:"getPrice"`
                Quantity    int     `json:"quantity"`
                Sector      string  `json:"sector"`
                UsdJpy      float64 `json:"usdJpy"`
                CurrentPrice float64 `json:"currentPrice"`
                PriceGets   float64 `json:"priceGets"`
                CurrentRate float64 `json:"currentRate"`
            } `json:"updateUsStock"`
        } `json:"data"`
    }
	t.Log("==========")
	t.Log(w.Body)

    err = json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v", err)
    }

    // レスポンスボディの内容の検証
    assert.Equal(t, "VTI", response.Data.UpdateUsStock.Code)
    assert.Equal(t, 234.0, response.Data.UpdateUsStock.GetPrice)
    assert.Equal(t, 15, response.Data.UpdateUsStock.Quantity)
    assert.Equal(t, "Index", response.Data.UpdateUsStock.Sector)
    assert.Equal(t, 135.0, response.Data.UpdateUsStock.UsdJpy)

    // データベースの更新内容を確認
    var updatedStock model.UsStock
    db.First(&updatedStock, "id = ?", createdUsStockID)

    assert.Equal(t, 234.0, updatedStock.GetPrice)
    assert.Equal(t, 15.0, updatedStock.Quantity)
    assert.Equal(t, 135.0, updatedStock.UsdJpy)
}


func TestDeleteUsStockE2E(t *testing.T) {
	db := test.SetupTestDB()
	router := graphql.SetupGraphQLServer(db, nil)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// テスト用データの追加
	usStock := model.UsStock{Code: "AAPL", UserId: 1, Quantity: 10, GetPrice: 100, Sector: "IT", UsdJpy: 133.9}
	db.Create(&usStock)

	// 作成されたレコードのIDを取得
	createdUsStockID := usStock.ID

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// createdUsStockIDを文字列に変換
	createdUsStockIDStr := strconv.FormatUint(uint64(createdUsStockID), 10)

	// GraphQLリクエストの実行
	query := fmt.Sprintf(`mutation {
		deleteUsStock(id: "%s")
	}`, createdUsStockIDStr)
	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			DeleteUsStock bool `json:"deleteUsStock"`
		} `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.True(t, response.Data.DeleteUsStock)

	// データベースから削除されたことを確認
	var stockAfterDelete model.UsStock
	result := db.First(&stockAfterDelete, "id = ?", createdUsStockID)
	assert.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
}
