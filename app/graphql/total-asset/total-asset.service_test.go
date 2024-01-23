package totalasset

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/lib/pq"
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

	mockStocks := []model.UsStock{
		{Code: "AAPL", GetPrice: 150, Quantity: 10, Sector: "Technology"},
	}
	mockStockRepo.On("FetchUsStockListById", mock.Anything, userId).Return(mockStocks, nil)

	mockMarketPrices := []marketPrice.MarketPriceDto{
		{Ticker: "AAPL", CurrentPrice: 155, PriceGets: 5, CurrentRate: 0.0333},
	}
	mockMarketPriceRepo.On("FetchMarketPriceList", mock.Anything, []string{"AAPL"}).Return(mockMarketPrices, nil)

	expectedUsdJpy := 133.69
    mockCurrencyRepo.On("FetchCurrentUsdJpy", mock.Anything).Return(expectedUsdJpy, nil)

	mockFunds := []model.JapanFund{
		{Code: "SP500", Name:"ｅＭＡＸＩＳ Ｓｌｉｍ 米国株式（Ｓ＆Ｐ５００）", GetPrice: 15523.81, GetPriceTotal: 761157.0,UserId: 1},
	}
	mockJapanFundRepo.On("FetchJapanFundListById", mock.Anything, userId).Return(mockFunds, nil)

	mockCryptos := []model.Crypto{
		{Code: "xrp", GetPrice: 88.0, Quantity: 2.0},
	}
	mockCryptoRepo.On("FetchCryptoListById", mock.Anything, userId).Return(mockCryptos, nil)

	mockMarketPrice := &marketCryptoRepo.Crypto{Name: "xrp", Price: 88.2}
	mockMarketCryptoRepo.On("FetchCryptoPrice", "xrp").Return(mockMarketPrice, nil)

	mockAssets := []model.FixedIncomeAsset{
		{Code: "Funds", UserId: 99, DividendRate: 3.5, GetPriceTotal: 100000.0, PaymentMonth: pq.Int64Array{6, 12}},
	}
	mockFixedIncomeAssetRepo.On("FetchFixedIncomeAssetListById", mock.Anything, userId).Return(mockAssets, nil)

	// テスト実行
	updateInput := generated.UpdateTotalAssetInput{ID: "1", CashJpy: 15000, CashUsd: 150}
	mockUpdatedAsset := &model.TotalAsset{CashJpy: 15000, CashUsd: 150, Stock:20000}
	mockTotalAssetRepo.On("UpdateTotalAsset", mock.Anything, mock.Anything).Return(mockUpdatedAsset, nil)

	// テスト対象メソッドの実行
	updatedAsset, err := service.UpdateTotalAsset(context.Background(), updateInput)
	assert.NoError(t, err)
	assert.NotNil(t, updatedAsset)
	assert.Equal(t, float64(15000), updatedAsset.CashJpy)
	assert.Equal(t, float64(150), updatedAsset.CashUsd)
	assert.Equal(t, float64(20000), updatedAsset.Stock)

	// モックの呼び出しを検証
	mockTotalAssetRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
