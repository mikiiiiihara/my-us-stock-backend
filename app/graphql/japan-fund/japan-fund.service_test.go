package japanfund

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	repo "my-us-stock-backend/app/repository/assets/fund"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUsStocks は UsStocks メソッドのテストです。
func TestJapanFundsService(t *testing.T) {
	mockRepo := repo.NewMockJapanFundRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewJapanFundService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockFunds := []model.JapanFund{
		{Code: "SP500", Name:"ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", GetPrice: 15523.81, GetPriceTotal: 761157.0,UserId: 1},
	}
	mockRepo.On("FetchJapanFundListById", mock.Anything, userId).Return(mockFunds, nil)

	// テスト対象メソッドの実行
	funds, err := service.JapanFunds(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, funds)
	assert.Len(t, funds, 1)

	assert.Equal(t, "0", funds[0].ID)
	assert.Equal(t, "SP500", funds[0].Code)
	assert.Equal(t, "ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", funds[0].Name)
	assert.Equal(t, 15523.81, funds[0].GetPrice)
	assert.Equal(t, 761157.0, funds[0].GetPriceTotal)
	assert.Equal(t, 25779.0, funds[0].CurrentPrice)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestCreateJapanFundService は TestCreateJapanFund メソッドのテストです。
func TestCreateJapanFundService(t *testing.T) {
	mockRepo := repo.NewMockJapanFundRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewJapanFundService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockAsset := &model.JapanFund{Code: "全世界株", Name:"ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）", GetPrice: 18609, GetPriceTotal: 400004,UserId: 1}

	input := repo.CreateJapanFundDto{
		Code: "全世界株", 
		Name:"ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）", 
		GetPrice: 18609, 
		GetPriceTotal: 400004.0,
		UserId: 1,
	}
	mockRepo.On("CreateJapanFund", mock.Anything, input).Return(mockAsset, nil)

	// テスト対象メソッドの実行
	serviceInput := generated.CreateJapanFundInput{
		Code: "全世界株", 
		Name:"ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）", 
		GetPrice: 18609, 
		GetPriceTotal: 400004.0,
	}
	newFund, err := service.CreateJapanFund(context.Background(), serviceInput)
	assert.NoError(t, err)
	assert.NotNil(t, newFund)

	assert.Equal(t, "0", newFund.ID)
	assert.Equal(t, "全世界株", newFund.Code)
	assert.Equal(t, "ｅＭＡＸＩＳ　Ｓｌｉｍ　全世界株式（除く日本）", newFund.Name)
	assert.Equal(t, 18609.0, newFund.GetPrice)
	assert.Equal(t, 400004.0, newFund.GetPriceTotal)
	assert.Equal(t, 22023.0, newFund.CurrentPrice)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestDeleteJapanFundService は DeleteJapanFund メソッドのテストです。
func TestDeleteJapanFundService(t *testing.T) {
	mockRepo := repo.NewMockJapanFundRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewJapanFundService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	deleteId := uint(1)
	mockRepo.On("DeleteJapanFund", mock.Anything, deleteId).Return(nil)

	// テスト対象メソッドの実行
	result, err := service.DeleteJapanFund(context.Background(), "1")
	assert.NoError(t, err)
	assert.True(t, result)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
