package crypto

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"github.com/stretchr/testify/mock"
)

type MockCryptoRepository struct {
	mock.Mock
}

func NewMockCryptoRepository() *MockCryptoRepository {
	return &MockCryptoRepository{}
}

func (m *MockCryptoRepository) FetchCryptoListById(ctx context.Context, userId uint) ([]model.Crypto, error) {
    args := m.Called(ctx, userId)
    return args.Get(0).([]model.Crypto), args.Error(1)
}

func (m *MockCryptoRepository) UpdateCrypto(ctx context.Context, dto UpdateCryptoDto) (*model.Crypto, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.Crypto), args.Error(1)
}

func (m *MockCryptoRepository) CreateCrypto(ctx context.Context, dto CreateCryptDto) (*model.Crypto, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.Crypto), args.Error(1)
}

func (m *MockCryptoRepository) DeleteCrypto(ctx context.Context, id uint) error{
	args := m.Called(ctx, id)
	return args.Error(0)
}