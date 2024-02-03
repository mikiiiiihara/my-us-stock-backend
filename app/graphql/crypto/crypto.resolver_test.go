package crypto

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCryptoService は MarketPriceService のモックです。
type MockCryptoService struct {
    mock.Mock
}

func (m *MockCryptoService) Cryptos(ctx context.Context) ([]*generated.Crypto, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*generated.Crypto), args.Error(1)
}

func (m *MockCryptoService) CreateCrypto(ctx context.Context, input generated.CreateCryptoInput) (*generated.Crypto, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.Crypto), args.Error(1)
}

func (m *MockCryptoService) UpdateCrypto(ctx context.Context, input generated.UpdateCryptoInput) (*generated.Crypto, error) {
    args := m.Called(ctx, input)
    return args.Get(0).(*generated.Crypto), args.Error(1)
}


func (m *MockCryptoService) DeleteCrypto(ctx context.Context, id string) (bool, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(bool), args.Error(1)
}

func TestCryptos(t *testing.T) {
    mockService := new(MockCryptoService)
    resolver := NewResolver(mockService)

    cryptos := []*generated.Crypto{
        {ID: "1",Code: "btc", GetPrice: 5047113.0, Quantity: 0.05, CurrentPrice: 5947113.84},
        {ID: "2",Code: "xrp",  GetPrice: 88.0, Quantity: 2,CurrentPrice: 88.2},
    }
    mockService.On("Cryptos", mock.Anything).Return(cryptos, nil)

    result, err := resolver.Cryptos(context.Background())
    
    assert.NoError(t, err)
    assert.Equal(t, cryptos, result)

    mockService.AssertExpectations(t)
}

func TestCreateCrypto(t *testing.T) {
    mockService := new(MockCryptoService)
    resolver := NewResolver(mockService)

	input := generated.CreateCryptoInput{
		Code: "xrp",
		GetPrice: 88.0, 
		Quantity: 2,
	}
	mockResponse := &generated.Crypto{
		ID: "1",
		Code: "xrp",
		GetPrice: 88.0, 
		Quantity: 2, 
		CurrentPrice: 88.2,
	}
    mockService.On("CreateCrypto", mock.Anything, input).Return(mockResponse, nil)

    result, err := resolver.CreateCrypto(context.Background(), input)
    
    assert.NoError(t, err)
    assert.Equal(t, mockResponse, result)

    mockService.AssertExpectations(t)
}

func TestUpdateCrypto(t *testing.T) {
    mockService := new(MockCryptoService)
    resolver := NewResolver(mockService)

    input := generated.UpdateCryptoInput{
        ID: "1",
        GetPrice: 403000.0, 
        Quantity: 1.2,
    }
    mockResponse := &generated.Crypto{
        ID: "1",
        Code: "eth",
        GetPrice: 403000.0, 
        Quantity: 1.2, 
        CurrentPrice: 413000.0,
    }
    mockService.On("UpdateCrypto", mock.Anything, input).Return(mockResponse, nil)

    result, err := resolver.UpdateCrypto(context.Background(), input)
    
    assert.NoError(t, err)
    assert.Equal(t, mockResponse, result)

    mockService.AssertExpectations(t)
}


func TestDeleteCrypto(t *testing.T) {
    mockService := new(MockCryptoService)
    resolver := NewResolver(mockService)

    mockService.On("DeleteCrypto", mock.Anything, "1").Return(true, nil)

    result, err := resolver.DeleteCrypto(context.Background(), "1")
    
    assert.NoError(t, err)
    assert.True(t, result)

    mockService.AssertExpectations(t)
}
