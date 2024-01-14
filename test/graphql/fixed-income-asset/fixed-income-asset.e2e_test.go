package fixedincomeasset

import (
	"encoding/json"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lib/pq"
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

func TestFixedIncomeAssets(t *testing.T) {
	db := test.SetupTestDB()
    router := graphql.SetupGraphQLServer(db, nil)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.FixedIncomeAsset{Code: "Funds", UserId: 1, DividendRate: 3.5, GetPriceTotal: 100000.0, PaymentMonth: pq.Int64Array{6, 12}})

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `query {
		fixedIncomeAssets{ id code getPriceTotal dividendRate usdJpy paymentMonth   }
	  }`
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            FixedIncomeAssets []struct {
                ID string `json:"id"`
				Code string `json:"code"`
				GetPriceTotal float64 `json:"getPriceTotal"`
				UsdJpy float64 `json:"usdJpy"`
                DividendRate float64 `json:"dividendRate"`
				PaymentMonth []int `json:"paymentMonth"`
            } `json:"fixedIncomeAssets"`
        } `json:"data"`
    }
	t.Log(w.Body)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
    if len(response.Data.FixedIncomeAssets) > 0 {
        assert.Equal(t, "Funds", response.Data.FixedIncomeAssets[0].Code)
		assert.Equal(t, 100000.0, response.Data.FixedIncomeAssets[0].GetPriceTotal)
		assert.Equal(t, 3.5, response.Data.FixedIncomeAssets[0].DividendRate)
		assert.Equal(t, 0.0, response.Data.FixedIncomeAssets[0].UsdJpy)
		assert.Equal(t, []int{6,12}, response.Data.FixedIncomeAssets[0].PaymentMonth)
    } else {
        t.Fatalf("Expected non-empty array")
    }
}

func TestCreateUsStockE2E(t *testing.T) {
	db := test.SetupTestDB()
    router := graphql.SetupGraphQLServer(db, nil)

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
		createFixedIncomeAsset(input: {
			code: "Funds-からだにユーグレナファンド",
			getPriceTotal: 110000.0,
			dividendRate: 1.8,
			paymentMonth: [3]
		}) {
			id
			code
			getPriceTotal
			dividendRate
			usdJpy
			paymentMonth
		}
	  }
	  `
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)
	  t.Log(w.Body)

	  t.Log(w.Body)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            CreateFixedIncomeAsset struct {
                ID string `json:"id"`
				Code string `json:"code"`
				GetPriceTotal float64 `json:"getPriceTotal"`
				UsdJpy float64 `json:"usdJpy"`
                DividendRate float64 `json:"dividendRate"`
				PaymentMonth []int `json:"paymentMonth"`
            } `json:"createFixedIncomeAsset"`
        } `json:"data"`
    }

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
	assert.Equal(t, "Funds-からだにユーグレナファンド", response.Data.CreateFixedIncomeAsset.Code)
	assert.Equal(t, 110000.0, response.Data.CreateFixedIncomeAsset.GetPriceTotal)
	assert.Equal(t, 1.8, response.Data.CreateFixedIncomeAsset.DividendRate)
	assert.Equal(t, 0.0, response.Data.CreateFixedIncomeAsset.UsdJpy)
	assert.Equal(t, []int{3}, response.Data.CreateFixedIncomeAsset.PaymentMonth)
}