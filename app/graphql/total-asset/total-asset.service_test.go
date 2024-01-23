package totalasset

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cryptoRepo "my-us-stock-backend/app/repository/assets/crypto"
	fixedIncomeAssetRepo "my-us-stock-backend/app/repository/assets/fixed-income"
	fundRepo "my-us-stock-backend/app/repository/assets/fund"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price"
	marketCryptoRepo "my-us-stock-backend/app/repository/market-price/crypto"
	"my-us-stock-backend/app/repository/market-price/currency"
	totalAssetRepo "my-us-stock-backend/app/repository/total-assets"
)

// MockTotalAssetRepository は TotalAssetRepository のモックです。
type MockTotalAssetRepository struct {
	mock.Mock
}

func (m *MockTotalAssetRepository) FetchTotalAssetListById(ctx context.Context, userId uint, day int) ([]model.TotalAsset, error) {
	args := m.Called(ctx, userId, day)
	return args.Get(0).([]model.TotalAsset), args.Error(1)
}

func (m *MockTotalAssetRepository) UpdateTotalAsset(ctx context.Context, dto totalAssetRepo.UpdateTotalAssetDto) (*model.TotalAsset, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.TotalAsset), args.Error(1)
}

func (m *MockTotalAssetRepository) CreateTodayTotalAsset(ctx context.Context, dto totalAssetRepo.CreateTotalAssetDto) (*model.TotalAsset, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.TotalAsset), args.Error(1)
}

func (m *MockTotalAssetRepository) FindTodayTotalAsset(ctx context.Context, userId uint) (*model.TotalAsset, error){
	args := m.Called(ctx, userId)
	return args.Get(0).(*model.TotalAsset), args.Error(1)
}

// TotalAssets メソッドのテスト
func TestTotalAssetsService(t *testing.T) {
	mockAuth := auth.NewMockAuthService()
	mockTotalAssetRepo := new(MockTotalAssetRepository)
	mockStockRepo := stock.NewMockUsStockRepository()
	mockMarketPriceRepo := marketPrice.NewMockMarketPriceRepository()
	mockCurrencyRepo := currency.NewMockCurrencyRepository()
	mockJapanFundRepo := fundRepo.NewMockJapanFundRepository()
	mockCryptoRepo := cryptoRepo.NewMockCryptoRepository()
	mockFixedIncomeAssetRepo := fixedIncomeAssetRepo.NewMockFixedIncomeAssetRepository()
	mockMarketCryptoRepo := marketCryptoRepo.NewMockCryptoRepository()
	service := NewTotalAssetService(mockAuth, mockTotalAssetRepo, mockStockRepo, mockMarketPriceRepo, mockCurrencyRepo, mockJapanFundRepo, mockCryptoRepo, mockFixedIncomeAssetRepo,mockMarketCryptoRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockAssets := []model.TotalAsset{
		{CashJpy: 10000, CashUsd: 100, Stock: 50000},
	}
	mockTotalAssetRepo.On("FetchTotalAssetListById", mock.Anything, userId, 30).Return(mockAssets, nil)

	// テスト対象メソッドの実行
	assets, err := service.TotalAssets(context.Background(), 30)
	assert.NoError(t, err)
	assert.Len(t, assets, 1)

	// モックの呼び出しを検証
	mockTotalAssetRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// UpdateTotalAsset メソッドのテスト
func TestUpdateTotalAssetService(t *testing.T) {
	mockAuth := auth.NewMockAuthService()
	mockTotalAssetRepo := new(MockTotalAssetRepository)
	mockStockRepo := stock.NewMockUsStockRepository()
	mockMarketPriceRepo := marketPrice.NewMockMarketPriceRepository()
	mockCurrencyRepo := currency.NewMockCurrencyRepository()
	mockJapanFundRepo := fundRepo.NewMockJapanFundRepository()
	mockCryptoRepo := cryptoRepo.NewMockCryptoRepository()
	mockFixedIncomeAssetRepo := fixedIncomeAssetRepo.NewMockFixedIncomeAssetRepository()
	mockMarketCryptoRepo := marketCryptoRepo.NewMockCryptoRepository()
	service := NewTotalAssetService(mockAuth, mockTotalAssetRepo, mockStockRepo, mockMarketPriceRepo, mockCurrencyRepo, mockJapanFundRepo, mockCryptoRepo, mockFixedIncomeAssetRepo,mockMarketCryptoRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	updateInput := generated.UpdateTotalAssetInput{ID: "1", CashJpy: 15000, CashUsd: 150}
	updateDto := totalAssetRepo.UpdateTotalAssetDto{ID: 1, CashJpy: &updateInput.CashJpy, CashUsd: &updateInput.CashUsd}
	mockUpdatedAsset := &model.TotalAsset{CashJpy: 15000, CashUsd: 150}
	mockTotalAssetRepo.On("UpdateTotalAsset", mock.Anything, updateDto).Return(mockUpdatedAsset, nil)

	// テスト対象メソッドの実行
	updatedAsset, err := service.UpdateTotalAsset(context.Background(), updateInput)
	assert.NoError(t, err)
	assert.NotNil(t, updatedAsset)
	assert.Equal(t, float64(15000), updatedAsset.CashJpy)
	assert.Equal(t, float64(150), updatedAsset.CashUsd)

	// モックの呼び出しを検証
	mockTotalAssetRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
