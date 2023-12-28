package stock

import (
	"context"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/repository/assets/stock"
	marketPrice "my-us-stock-backend/app/repository/market-price"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsStockRepository は UsStockRepository のモックです。
type MockUsStockRepository struct {
	mock.Mock
}

func (m *MockUsStockRepository) FetchUsStockListById(ctx context.Context, userId uint) ([]model.UsStock, error) {
    args := m.Called(ctx, userId)
    return args.Get(0).([]model.UsStock), args.Error(1)
}

func (m *MockUsStockRepository) UpdateUsStock(ctx context.Context, dto stock.UpdateUsStockDto) (*model.UsStock, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.UsStock), args.Error(1)
}

func (m *MockUsStockRepository) CreateUsStock(ctx context.Context, dto stock.CreateUsStockDto) (*model.UsStock, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.UsStock), args.Error(1)
}

func (m *MockUsStockRepository) DeleteUsStock(ctx context.Context, id uint) error{
	args := m.Called(ctx, id)
	return args.Error(1)
}

// MockMarketPriceRepository は MarketPriceRepository のモックです。
type MockMarketPriceRepository struct {
	mock.Mock
}

func (m *MockMarketPriceRepository) FetchMarketPriceList(ctx context.Context, tickers []string) ([]marketPrice.MarketPriceDto, error) {
	args := m.Called(ctx, tickers)
	return args.Get(0).([]marketPrice.MarketPriceDto), args.Error(1)
}

func (m *MockMarketPriceRepository) FetchDividend(ctx context.Context, ticker string) (*marketPrice.DividendEntity, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*marketPrice.DividendEntity), args.Error(1)
}

// MockAuthService は AuthService のモックです。
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) FetchUserIdAccessToken(ctx context.Context) (uint, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockAuthService) RefreshAccessToken(c *gin.Context) (string, error) {
    args := m.Called(c)
    return args.String(0), args.Error(1)
}

func (m *MockAuthService) SignIn(ctx context.Context, c *gin.Context) (*model.User, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) SignUp(ctx context.Context, c *gin.Context) (*model.User, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) SendAuthResponse(ctx context.Context, c *gin.Context, user *model.User, code int) {
    m.Called(ctx, c, user, code)
}

// TestUsStocks は UsStocks メソッドのテストです。
func TestUsStocksService(t *testing.T) {
	mockStockRepo := new(MockUsStockRepository)
	mockMarketPriceRepo := new(MockMarketPriceRepository)
	mockAuth := new(MockAuthService)
	service := NewUsStockService(mockStockRepo, mockAuth, mockMarketPriceRepo)

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

	mockDividend := &marketPrice.DividendEntity{DividendTotal: 1.5}
	mockMarketPriceRepo.On("FetchDividend", mock.Anything, "AAPL").Return(mockDividend, nil)

	// テスト対象メソッドの実行
	usStocks, err := service.UsStocks(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, usStocks)
	assert.Len(t, usStocks, 1)

	assert.Equal(t, "0", usStocks[0].ID)
	assert.Equal(t, "AAPL", usStocks[0].Code)
	assert.Equal(t, 150.0, usStocks[0].GetPrice)
	assert.Equal(t, 10.0, usStocks[0].Quantity)
	assert.Equal(t, "Technology", usStocks[0].Sector)
	assert.Equal(t, 1.5, usStocks[0].Dividend)
	assert.Equal(t, 155.0, usStocks[0].CurrentPrice)
	assert.Equal(t, 5.0, usStocks[0].PriceGets)
	assert.Equal(t, 0.0333, usStocks[0].CurrentRate)

	// モックの呼び出しを検証
	mockStockRepo.AssertExpectations(t)
	mockMarketPriceRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// TestCreateUsStockService は TestCreateUsStock メソッドのテストです。
func TestCreateUsStockService(t *testing.T) {
	mockStockRepo := new(MockUsStockRepository)
	mockMarketPriceRepo := new(MockMarketPriceRepository)
	mockAuth := new(MockAuthService)
	service := NewUsStockService(mockStockRepo, mockAuth, mockMarketPriceRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockStock := &model.UsStock{Code: "AAPL", GetPrice: 150, Quantity: 10, Sector: "Technology",UsdJpy: 133.0}

	input := stock.CreateUsStockDto{
		Code: "AAPL",
		GetPrice: 150, 
		Quantity: 10, 
		Sector: "Technology",
		UsdJpy: 133.0,
		UserId: 1,
	}
	mockStockRepo.On("CreateUsStock", mock.Anything, input).Return(mockStock, nil)

	mockMarketPrices := []marketPrice.MarketPriceDto{
		{Ticker: "AAPL", CurrentPrice: 155, PriceGets: 5, CurrentRate: 0.0333},
	}
	mockMarketPriceRepo.On("FetchMarketPriceList", mock.Anything, []string{"AAPL"}).Return(mockMarketPrices, nil)

	mockDividend := &marketPrice.DividendEntity{DividendTotal: 1.5}
	mockMarketPriceRepo.On("FetchDividend", mock.Anything, "AAPL").Return(mockDividend, nil)

	// テスト対象メソッドの実行
	serviceInput := generated.CreateUsStockInput{
		Code: "AAPL",
		GetPrice: 150, 
		Quantity: 10, 
		Sector: "Technology",
		UsdJpy: 133.0,
	}
	usStock, err := service.CreateUsStock(context.Background(), serviceInput)
	assert.NoError(t, err)
	assert.NotNil(t, usStock)

	assert.Equal(t, "0", usStock.ID)
	assert.Equal(t, "AAPL", usStock.Code)
	assert.Equal(t, 150.0, usStock.GetPrice)
	assert.Equal(t, 10.0, usStock.Quantity)
	assert.Equal(t, "Technology", usStock.Sector)
	assert.Equal(t, 1.5, usStock.Dividend)
	assert.Equal(t, 155.0, usStock.CurrentPrice)
	assert.Equal(t, 5.0, usStock.PriceGets)
	assert.Equal(t, 0.0333, usStock.CurrentRate)

	// モックの呼び出しを検証
	mockStockRepo.AssertExpectations(t)
	mockMarketPriceRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
