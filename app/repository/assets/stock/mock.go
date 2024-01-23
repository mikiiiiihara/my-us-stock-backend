package stock

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"github.com/stretchr/testify/mock"
)

// MockUsStockRepository は UsStockRepository のモックです。
type MockUsStockRepository struct {
	mock.Mock
}

func NewMockUsStockRepository() *MockUsStockRepository {
	return &MockUsStockRepository{}
}

func (m *MockUsStockRepository) FetchUsStockListById(ctx context.Context, userId uint) ([]model.UsStock, error) {
    args := m.Called(ctx, userId)
    return args.Get(0).([]model.UsStock), args.Error(1)
}

func (m *MockUsStockRepository) UpdateUsStock(ctx context.Context, dto UpdateUsStockDto) (*model.UsStock, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.UsStock), args.Error(1)
}

func (m *MockUsStockRepository) CreateUsStock(ctx context.Context, dto CreateUsStockDto) (*model.UsStock, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.UsStock), args.Error(1)
}

func (m *MockUsStockRepository) DeleteUsStock(ctx context.Context, id uint) error{
	args := m.Called(ctx, id)
	return args.Error(0)
}