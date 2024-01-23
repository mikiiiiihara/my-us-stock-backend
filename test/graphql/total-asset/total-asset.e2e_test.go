package totalasset

import (
	"encoding/json"
	"fmt"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/test"
	"my-us-stock-backend/test/graphql"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTotalAssetsE2E(t *testing.T) {
	db := test.SetupTestDB()
	router := graphql.SetupGraphQLServer(db, nil)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// テスト用データの追加
	db.Create(&model.TotalAsset{UserId: 1, CashJpy: 10000, CashUsd: 100, Stock: 50000})

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// GraphQLリクエストの実行
	query := `query {
		totalAssets(day: 30) {
			id cashJpy cashUsd stock fund crypto fixedIncomeAsset createdAt
		}
	}`
	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			TotalAssets []struct {
				ID               string  `json:"id"`
				CashJpy          float64 `json:"cashJpy"`
				CashUsd          float64 `json:"cashUsd"`
				Stock            float64 `json:"stock"`
				Fund             float64 `json:"fund"`
				Crypto           float64 `json:"crypto"`
				FixedIncomeAsset float64 `json:"fixedIncomeAsset"`
				CreatedAt        string  `json:"createdAt"`
			} `json:"totalAssets"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	// レスポンスボディの内容の検証
	assert.NotEmpty(t, response.Data.TotalAssets)
	for _, asset := range response.Data.TotalAssets {
		assert.Greater(t, asset.CashJpy, 0.0)
		assert.Greater(t, asset.CashUsd, 0.0)
		assert.Greater(t, asset.Stock, 0.0)
	}
}

func TestUpdateTotalAssetE2E(t *testing.T) {
	db := test.SetupTestDB()
	router := graphql.SetupGraphQLServer(db, nil)

	// テスト用HTTPサーバーのセットアップ
	ts := httptest.NewServer(router)
	defer ts.Close()

	// テスト用データの追加
	totalAsset := model.TotalAsset{UserId: 1, CashJpy: 10000, CashUsd: 100, Stock: 50000}
	db.Create(&totalAsset)

	// 作成されたレコードのIDを取得
	createdTotalAssetID := totalAsset.ID

	// ダミーのアクセストークンを生成
	token, err := graphql.GenerateTestAccessTokenForUserId(1)
	if err != nil {
		t.Fatalf("Failed to generate test access token: %v", err)
	}

	// createdTotalAssetIDを文字列に変換
	createdTotalAssetIDStr := strconv.FormatUint(uint64(createdTotalAssetID), 10)

	// GraphQLリクエストの実行
	query := fmt.Sprintf(`mutation {
		updateTotalAsset(input: {
			id: "%s"
			cashJpy: 15000
			cashUsd: 150
		}) {
			id cashJpy cashUsd stock fund crypto fixedIncomeAsset createdAt
		}
	}`, createdTotalAssetIDStr)
	w := graphql.ExecuteGraphQLRequestWithToken(ts.URL, query, token)

	// レスポンスボディの解析
	var response struct {
		Data struct {
			UpdateTotalAsset struct {
				ID               string  `json:"id"`
				CashJpy          float64 `json:"cashJpy"`
				CashUsd          float64 `json:"cashUsd"`
				Stock            float64 `json:"stock"`
				Fund             float64 `json:"fund"`
				Crypto           float64 `json:"crypto"`
				FixedIncomeAsset float64 `json:"fixedIncomeAsset"`
				CreatedAt        string  `json:"createdAt"`
			} `json:"updateTotalAsset"`
		} `json:"data"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// レスポンスボディの内容の検証
	assert.Equal(t, createdTotalAssetIDStr, response.Data.UpdateTotalAsset.ID)
	assert.Equal(t, 15000.0, response.Data.UpdateTotalAsset.CashJpy)
	assert.Equal(t, 150.0, response.Data.UpdateTotalAsset.CashUsd)
}
