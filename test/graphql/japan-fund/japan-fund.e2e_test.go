package japanfund

import (
	"encoding/json"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http"
	"net/http/httptest"
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

func TestJapanFunds(t *testing.T) {
	db := test.SetupTestDB()
    router := graphql.SetupGraphQLServer(db, nil)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.JapanFund{Code: "SP500", Name:"ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", GetPrice: 15523.81, GetPriceTotal: 761157.0,UserId: 1})

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `query {
		japanFunds{ id code name getPrice getPriceTotal currentPrice   }
	  }`
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            JapanFunds []struct {
                ID string `json:"id"`
				Code string `json:"code"`
				Name string `json:"name"`
				GetPrice float64 `json:"getPrice"`
				GetPriceTotal float64 `json:"getPriceTotal"`
				CurrentPrice float64 `json:"currentPrice"`
            } `json:"japanFunds"`
        } `json:"data"`
    }
	t.Log(w.Body)

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
    if len(response.Data.JapanFunds) > 0 {
        assert.Equal(t, "SP500", response.Data.JapanFunds[0].Code)
		assert.Equal(t, "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", response.Data.JapanFunds[0].Name)
		assert.Equal(t, 15523.81, response.Data.JapanFunds[0].GetPrice)
		assert.Equal(t, 761157.0, response.Data.JapanFunds[0].GetPriceTotal)
		assert.Equal(t, 25369.0, response.Data.JapanFunds[0].CurrentPrice)
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
		createJapanFund(input: {
			code: "全世界株",
			name: "ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）",
			getPrice: 18609.0,
			getPriceTotal: 400004.0
		}) {
			id 
			code 
			name 
			getPrice 
			getPriceTotal 
			currentPrice
		}
	  }
	  `
	  w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)
	  t.Log(w.Body)

	  t.Log(w.Body)

	// レスポンスボディの解析
    var response struct {
        Data struct {
            CreateJapanFund struct {
                ID string `json:"id"`
				Code string `json:"code"`
				Name string `json:"name"`
				GetPrice float64 `json:"getPrice"`
				GetPriceTotal float64 `json:"getPriceTotal"`
				CurrentPrice float64 `json:"currentPrice"`
            } `json:"createJapanFund"`
        } `json:"data"`
    }

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

    // レスポンスボディの内容の検証
	assert.Equal(t, "全世界株", response.Data.CreateJapanFund.Code)
	assert.Equal(t, "ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）", response.Data.CreateJapanFund.Name)
	assert.Equal(t, 18609.0, response.Data.CreateJapanFund.GetPrice)
	assert.Equal(t, 400004.0, response.Data.CreateJapanFund.GetPriceTotal)
	assert.Equal(t, 21682.0, response.Data.CreateJapanFund.CurrentPrice)
}