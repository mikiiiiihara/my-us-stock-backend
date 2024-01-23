package fixedincome

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"github.com/stretchr/testify/mock"
)

type MockFixedIncomeAssetRepository struct {
	mock.Mock
}

func NewMockFixedIncomeAssetRepository() *MockFixedIncomeAssetRepository {
	return &MockFixedIncomeAssetRepository{}
}

func (m *MockFixedIncomeAssetRepository) FetchFixedIncomeAssetListById(ctx context.Context, userId uint) ([]model.FixedIncomeAsset, error) {
    args := m.Called(ctx, userId)
    return args.Get(0).([]model.FixedIncomeAsset), args.Error(1)
}

func (m *MockFixedIncomeAssetRepository) UpdateFixedIncomeAsset(ctx context.Context, dto UpdateFixedIncomeDto) (*model.FixedIncomeAsset, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.FixedIncomeAsset), args.Error(1)
}

func (m *MockFixedIncomeAssetRepository)CreateFixedIncomeAsset(ctx context.Context, dto CreateFixedIncomeDto) (*model.FixedIncomeAsset, error){
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.FixedIncomeAsset), args.Error(1)
}

func (m *MockFixedIncomeAssetRepository) DeleteFixedIncomeAsset(ctx context.Context, id uint) error{
	args := m.Called(ctx, id)
	return args.Error(0)
}
