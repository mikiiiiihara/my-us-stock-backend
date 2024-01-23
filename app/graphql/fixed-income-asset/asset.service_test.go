package fixedincomeasset

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	repo "my-us-stock-backend/app/repository/assets/fixed-income"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestUsStocks は UsStocks メソッドのテストです。
func TestFixedIncomeAssetsService(t *testing.T) {
	mockRepo := repo.NewMockFixedIncomeAssetRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewAssetService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockAssets := []model.FixedIncomeAsset{
		{Code: "Funds", UserId: 99, DividendRate: 3.5, GetPriceTotal: 100000.0, PaymentMonth: pq.Int64Array{6, 12}},
	}
	mockRepo.On("FetchFixedIncomeAssetListById", mock.Anything, userId).Return(mockAssets, nil)

	// テスト対象メソッドの実行
	assets, err := service.FixedIncomeAssets(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, assets)
	assert.Len(t, assets, 1)

	assert.Equal(t, "0", assets[0].ID)
	assert.Equal(t, "Funds", assets[0].Code)
	assert.Equal(t, 100000.0, assets[0].GetPriceTotal)
	assert.Equal(t, 3.5, assets[0].DividendRate)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestCreateUsStockService は TestCreateUsStock メソッドのテストです。
func TestCreateFixedIncomeAssetService(t *testing.T) {
	mockRepo := repo.NewMockFixedIncomeAssetRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewAssetService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockAsset := &model.FixedIncomeAsset{Code:"Bankers",DividendRate: 3.5, GetPriceTotal: 100000.0,PaymentMonth: pq.Int64Array{6, 12},UserId: 1,}

	paymentMonth := []int64{6,12}

	input := repo.CreateFixedIncomeDto{
        Code:   "Bankers",
		DividendRate: 3.5, 
		GetPriceTotal: 100000.0,
        PaymentMonth: paymentMonth,
		UserId: 1,
	}
	mockRepo.On("CreateFixedIncomeAsset", mock.Anything, input).Return(mockAsset, nil)

	// テスト対象メソッドの実行
	serviceInput := generated.CreateFixedIncomeAssetInput{
        Code:   "Bankers",
		DividendRate: 3.5, 
		GetPriceTotal: 100000.0,
        PaymentMonth: []int{6,12},
	}
	newAsset, err := service.CreateFixedIncomeAsset(context.Background(), serviceInput)
	assert.NoError(t, err)
	assert.NotNil(t, newAsset)

	assert.Equal(t, "0", newAsset.ID)
	assert.Equal(t, "Bankers", newAsset.Code)
	assert.Equal(t, 100000.0, newAsset.GetPriceTotal)
	assert.Equal(t, 3.5, newAsset.DividendRate)
	assert.Equal(t, 3.5, newAsset.DividendRate)
	assert.Equal(t, []int{6,12}, newAsset.PaymentMonth)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestDeleteFixedIncomeAssetService は DeleteFixedIncomeAsset メソッドのテストです。
func TestDeleteFixedIncomeAssetService(t *testing.T) {
	mockRepo := repo.NewMockFixedIncomeAssetRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewAssetService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	deleteId := uint(1)
	mockRepo.On("DeleteFixedIncomeAsset", mock.Anything, deleteId).Return(nil)

	// テスト対象メソッドの実行
	result, err := service.DeleteFixedIncomeAsset(context.Background(), "1")
	assert.NoError(t, err)
	assert.True(t, result)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
