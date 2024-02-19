package fund

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"github.com/stretchr/testify/mock"
)

// MockFundPriceRepository is a mock of the FundPriceRepository interface
type MockFundPriceRepository struct {
	mock.Mock
}

// NewMockCryptoRepository は新しい NewMockCryptoRepository を作成し、初期設定を行います。
func NewMockFundPriceRepository() *MockFundPriceRepository {
	return &MockFundPriceRepository{}
}

func (m *MockFundPriceRepository) FetchFundPriceList(ctx context.Context) ([]model.FundPrice, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.FundPrice), args.Error(1)
}

func (m *MockFundPriceRepository) FindFundPriceByCode(ctx context.Context, code string) (*model.FundPrice, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*model.FundPrice), args.Error(1)
}

func (m *MockFundPriceRepository) UpdateFundPrice(ctx context.Context, dto UpdateFundPriceDto) (*model.FundPrice, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.FundPrice), args.Error(1)
}

func (m *MockFundPriceRepository) CreateFundPrice(ctx context.Context, dto CreateFundPriceDto) (*model.FundPrice, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.FundPrice), args.Error(1)
}