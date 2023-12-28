package marketprice

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

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

func TestFetchMarketPriceList(t *testing.T) {
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

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	tickers := []string{"AAPL", "KO"}
	prices, err := repo.FetchMarketPriceList(context.Background(), tickers)

	// アサーション
	assert.NoError(t, err)
	assert.Len(t, prices, 2)
	// AAPL
	assert.Equal(t, "AAPL", prices[0].Ticker)
	assert.Equal(t, 189.84, prices[0].CurrentPrice)
	assert.Equal(t, 0.0685, prices[0].CurrentRate)
	assert.Equal(t, 0.13, prices[0].PriceGets)
	// KO
	assert.Equal(t, "KO", prices[1].Ticker)
	assert.Equal(t, 57.205, prices[1].CurrentPrice)
	assert.Equal(t, 0.0685, prices[0].CurrentRate)
	assert.Equal(t, 0.13, prices[0].PriceGets)
}

func TestFetchMarketPriceList_NoData(t *testing.T) {
	// モックの HTTP レスポンスを設定（空のリスト）
	mockResponseBody := `[]`

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	tickers := []string{"AAPLKK"}
	_, err := repo.FetchMarketPriceList(context.Background(), tickers)

	// アサーション：期待されるエラーメッセージをチェック
	expectedErrorMessage := "the specified tickers were not found" // 期待されるエラーメッセージに更新
	assert.Equal(t, expectedErrorMessage, err.Error())
}

// 指定したティッカーの現在の配当情報を取得する
func TestFetchDividend(t *testing.T) {
	// モックの HTTP レスポンスを設定
	mockResponseBody := `{
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

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	dividend, err := repo.FetchDividend(context.Background(), "AAPL")

	// アサーション
	assert.NoError(t, err)
	// AAPL
	assert.Equal(t, "AAPL", dividend.Ticker)
	assert.Equal(t, 0.238, dividend.Dividend)
	assert.Equal(t, []int{2, 5, 8, 11}, dividend.DividendFixedMonth)
    assert.Equal(t, []int{2, 5, 8, 11}, dividend.DividendMonth)
    assert.Equal(t, 4, dividend.DividendTime)
    assert.Equal(t, 0.95, dividend.DividendTotal)
}

// 配当支払月に重複があった場合、重複なして出力される
func TestFetchDividend_PaymentDate_Duplicated(t *testing.T) {
	// モックの HTTP レスポンスを設定
	mockResponseBody := `{
		"symbol": "AAPL",
		"historical": [
			{
				"date": "2023-11-10",
				"label": "November 10, 23",
				"adjDividend": 0.24,
				"dividend": 0.24,
				"recordDate": "2023-11-13",
				"paymentDate": "2023-08-17",
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
				"paymentDate": "2023-08-17",
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

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	dividend, err := repo.FetchDividend(context.Background(), "AAPL")

	// アサーション
	assert.NoError(t, err)
	// AAPL
	assert.Equal(t, "AAPL", dividend.Ticker)
	assert.Equal(t, 0.238, dividend.Dividend)
	assert.Equal(t, []int{2, 5, 8, 11}, dividend.DividendFixedMonth)
    assert.Equal(t, []int{2, 8}, dividend.DividendMonth)
    assert.Equal(t, 4, dividend.DividendTime)
    assert.Equal(t, 0.95, dividend.DividendTotal)
}

// 配当付与月が取得できない場合、配当付与月=配当権利落月とする。
func TestFetchDividend_No_PaymentDate(t *testing.T) {
	// モックの HTTP レスポンスを設定
	mockResponseBody := `{
		"symbol": "AAPL",
		"historical": [
			{
				"date": "2023-11-10",
				"label": "November 10, 23",
				"adjDividend": 0.24,
				"dividend": 0.24,
				"recordDate": "2023-11-13",
				"paymentDate": "",
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
				"paymentDate": "2023-05-15",
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

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	dividend, err := repo.FetchDividend(context.Background(), "AAPL")

	// アサーション
	assert.NoError(t, err)
	// AAPL
	assert.Equal(t, "AAPL", dividend.Ticker)
	assert.Equal(t, 0.237, dividend.Dividend)
	assert.Equal(t, []int{2, 5, 8}, dividend.DividendFixedMonth)
    assert.Equal(t, []int{2, 5, 8}, dividend.DividendMonth)
    assert.Equal(t, 3, dividend.DividendTime)
    assert.Equal(t, 0.71, dividend.DividendTotal)
}

// 指定したティッカーの現在の配当情報を取得する(配当頻度が0回)
func TestFetchDividend_No_Dividend(t *testing.T) {
	// モックの HTTP レスポンスを設定
	mockResponseBody := `{
		"symbol": "TSLA",
		"historical": []
	}
	`

	// モックの HTTP クライアントを設定
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

	// リポジトリを初期化
	repo := NewMarketPriceRepository(mockHTTPClient)

	// テストの実行
	dividend, err := repo.FetchDividend(context.Background(), "TSLA")

	// アサーション
	assert.NoError(t, err)
	// AAPL
	assert.Equal(t, "TSLA", dividend.Ticker)
	assert.Equal(t, 0.0, dividend.Dividend)
	assert.Equal(t, []int(nil), dividend.DividendFixedMonth)
    assert.Equal(t, []int(nil), dividend.DividendMonth)
    assert.Equal(t, 0, dividend.DividendTime)
    assert.Equal(t, 0.0, dividend.DividendTotal)
}