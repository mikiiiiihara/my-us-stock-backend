package japanfund

import (
	"encoding/json"
	"fmt"
	"my-us-stock-backend/app/database/model"
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

func TestJapanFunds(t *testing.T) {
	db := test.SetupTestDB()
    router := graphql.SetupGraphQLServer(db, nil)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用のユーザーを作成
    db.Create(&model.JapanFund{Code: "SP500", Name:"ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", GetPrice: 15523.81, GetPriceTotal: 761157.0,UserId: 1})
    expectedFundPrice := model.FundPrice{Name: "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", Code: "SP500", Price: 27000.0}
    db.Create(&expectedFundPrice)

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
		assert.Equal(t, expectedFundPrice.Price, response.Data.JapanFunds[0].CurrentPrice)
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

	expectedFundPrice := model.FundPrice{Name: "ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）", Code: "全世界株", Price: 23000.0}
    db.Create(&expectedFundPrice)

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
	assert.Equal(t, expectedFundPrice.Price, response.Data.CreateJapanFund.CurrentPrice)
}

func TestUpdateJapanFundE2E(t *testing.T) {
    db := test.SetupTestDB()
    router := graphql.SetupGraphQLServer(db, nil)

    // テスト用HTTPサーバーのセットアップ
    ts := httptest.NewServer(router)
    defer ts.Close()

    // テスト用データの追加
    japanFund := model.JapanFund{Code: "JPX400", Name: "日経400", GetPrice: 12345.67, GetPriceTotal: 234567.89, UserId: 1}
    db.Create(&japanFund)

	expectedFundPrice := model.FundPrice{Name: "日経400", Code: "JPX400", Price: 18000.0}
    db.Create(&expectedFundPrice)

    // 作成されたレコードのIDを取得
    createdJapanFundID := japanFund.ID

    // ダミーのアクセストークンを生成
    token, err := graphql.GenerateTestAccessTokenForUserId(1)
    if err != nil {
        t.Fatalf("Failed to generate test access token: %v", err)
    }

    // GraphQLリクエストの実行
    updateQuery := fmt.Sprintf(`mutation {
        updateJapanFund(input: {
            id: "%s",
            getPrice: 13000.0,
            getPriceTotal: 260000.0
        }) {
            id
            code
            name
            getPrice
            getPriceTotal
            currentPrice
        }
    }`, strconv.FormatUint(uint64(createdJapanFundID), 10))

    w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, updateQuery, token)

    // レスポンスボディの解析
    var response struct {
        Data struct {
            UpdateJapanFund struct {
                ID            string  `json:"id"`
                Code          string  `json:"code"`
                Name          string  `json:"name"`
                GetPrice      float64 `json:"getPrice"`
                GetPriceTotal float64 `json:"getPriceTotal"`
                CurrentPrice  float64 `json:"currentPrice"`
            } `json:"updateJapanFund"`
        } `json:"data"`
    }

    err = json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Fatalf("Failed to parse response body: %v", err)
    }

    // レスポンスボディの内容の検証
    assert.Equal(t, strconv.FormatUint(uint64(createdJapanFundID), 10), response.Data.UpdateJapanFund.ID)
    assert.Equal(t, 13000.0, response.Data.UpdateJapanFund.GetPrice)
    assert.Equal(t, 260000.0, response.Data.UpdateJapanFund.GetPriceTotal)
	assert.Equal(t, expectedFundPrice.Price, response.Data.UpdateJapanFund.CurrentPrice)

    // データベースの更新内容を確認
    var updatedFund model.JapanFund
    result := db.First(&updatedFund, "id = ?", createdJapanFundID)
    assert.NoError(t, result.Error)
    assert.Equal(t, 13000.0, updatedFund.GetPrice)
    assert.Equal(t, 260000.0, updatedFund.GetPriceTotal)
}


func TestDeleteJapanFundE2E(t *testing.T) {
	db := test.SetupTestDB()
	router := graphql.SetupGraphQLServer(db, nil)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// テスト用データの追加
	japanFund := model.JapanFund{Code: "JPX400", Name: "日経400", GetPrice: 12345.67, GetPriceTotal: 234567.89, UserId: 1}
	db.Create(&japanFund)

	// 作成されたレコードのIDを取得
	createdJapanFundID := japanFund.ID

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// createdJapanFundIDを文字列に変換
	createdJapanFundIDStr := strconv.FormatUint(uint64(createdJapanFundID), 10)

	// GraphQLリクエストの実行
	query := fmt.Sprintf(`mutation {
		deleteJapanFund(id: "%s")
	}`, createdJapanFundIDStr)
	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			DeleteJapanFund bool `json:"deleteJapanFund"`
		} `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.True(t, response.Data.DeleteJapanFund)

	// データベースから削除されたことを確認
	var fundAfterDelete model.JapanFund
	result := db.First(&fundAfterDelete, "id = ?", createdJapanFundID)
	assert.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
}
