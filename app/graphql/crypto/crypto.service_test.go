package crypto

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"my-us-stock-backend/app/repository/assets/crypto"
	marketPrice "my-us-stock-backend/app/repository/market-price/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCryptoRepository struct {
	mock.Mock
}

func (m *MockCryptoRepository) FetchCryptoListById(ctx context.Context, userId uint) ([]model.Crypto, error) {
    args := m.Called(ctx, userId)
    return args.Get(0).([]model.Crypto), args.Error(1)
}

func (m *MockCryptoRepository) UpdateCrypto(ctx context.Context, dto crypto.UpdateCryptoDto) (*model.Crypto, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.Crypto), args.Error(1)
}

func (m *MockCryptoRepository) CreateCrypto(ctx context.Context, dto crypto.CreateCryptDto) (*model.Crypto, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.Crypto), args.Error(1)
}

func (m *MockCryptoRepository) DeleteCrypto(ctx context.Context, id uint) error{
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCryptosService(t *testing.T) {
	mockCryptoRepo := new(MockCryptoRepository)
    mockMarkeCryptoRepo := marketPrice.NewMockCryptoRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewCryptoService(mockCryptoRepo, mockAuth, mockMarkeCryptoRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockCryptos := []model.Crypto{
		{Code: "xrp", GetPrice: 88.0, Quantity: 2.0},
	}
	mockCryptoRepo.On("FetchCryptoListById", mock.Anything, userId).Return(mockCryptos, nil)


	mockMarketPrice := &marketPrice.Crypto{Name: "xrp", Price: 88.2}
	mockMarkeCryptoRepo.On("FetchCryptoPrice", "xrp").Return(mockMarketPrice, nil)

	// テスト対象メソッドの実行
	cryptos, err := service.Cryptos(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, cryptos)
	assert.Len(t, cryptos, 1)

	assert.Equal(t, "xrp", cryptos[0].Code)
	assert.Equal(t, 88.0, cryptos[0].GetPrice)
	assert.Equal(t, 2.0, cryptos[0].Quantity)
	assert.Equal(t, 88.2, cryptos[0].CurrentPrice)

	// モックの呼び出しを検証
	mockCryptoRepo.AssertExpectations(t)
	mockMarkeCryptoRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestCreateCryptoService(t *testing.T) {
	mockCryptoRepo := new(MockCryptoRepository)
    mockMarketCryptoRepo := marketPrice.NewMockCryptoRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewCryptoService(mockCryptoRepo, mockAuth, mockMarketCryptoRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockCrypto := &model.Crypto{Code: "btc", GetPrice: 5047113.0, Quantity: 0.05}

	input := crypto.CreateCryptDto{
		Code: "btc",
		GetPrice: 5047113.0, 
		Quantity: 0.05, 
		UserId: 1,
	}
	mockCryptoRepo.On("CreateCrypto", mock.Anything, input).Return(mockCrypto, nil)

	mockMarketPrice := &marketPrice.Crypto{Name: "btc", Price: 5947113.2}
	mockMarketCryptoRepo.On("FetchCryptoPrice", "btc").Return(mockMarketPrice, nil)

	// テスト対象メソッドの実行
	serviceInput := generated.CreateCryptoInput{
		Code: "btc",
		GetPrice: 5047113.0, 
		Quantity: 0.05, 
	}
	usStock, err := service.CreateCrypto(context.Background(), serviceInput)
	assert.NoError(t, err)
	assert.NotNil(t, usStock)

	assert.Equal(t, "btc", usStock.Code)
	assert.Equal(t, 5047113.0, usStock.GetPrice)
	assert.Equal(t, 0.05, usStock.Quantity)
	assert.Equal(t, 5947113.2, usStock.CurrentPrice)

	// モックの呼び出しを検証
	mockCryptoRepo.AssertExpectations(t)
	mockMarketCryptoRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestDeleteCryptoService(t *testing.T) {
	mockCryptoRepo := new(MockCryptoRepository)
	mockMarketCryptoRepo := marketPrice.NewMockCryptoRepository()
	mockAuth := auth.NewMockAuthService()
	service := NewCryptoService(mockCryptoRepo, mockAuth, mockMarketCryptoRepo)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	deleteId := uint(1)
	mockCryptoRepo.On("DeleteCrypto", mock.Anything, deleteId).Return(nil)

	// テスト対象メソッドの実行
	result, err := service.DeleteCrypto(context.Background(), "1")
	assert.NoError(t, err)
	assert.True(t, result)

	// モックの呼び出しを検証
	mockCryptoRepo.AssertExpectations(t)
	mockMarketCryptoRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
