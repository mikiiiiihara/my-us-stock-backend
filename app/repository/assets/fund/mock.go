package fund

import (
	"context"
	"my-us-stock-backend/app/database/model"

	"github.com/stretchr/testify/mock"
)

type MockJapanFundRepository struct {
	mock.Mock
}

func NewMockJapanFundRepository() *MockJapanFundRepository {
	return &MockJapanFundRepository{}
}

func (m *MockJapanFundRepository) FetchJapanFundListById(ctx context.Context, userId uint) ([]model.JapanFund, error) {
    args := m.Called(ctx, userId)
    return args.Get(0).([]model.JapanFund), args.Error(1)
}

func (m *MockJapanFundRepository) UpdateJapanFund(ctx context.Context, dto UpdateJapanFundDto) (*model.JapanFund, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.JapanFund), args.Error(1)
}

func (m *MockJapanFundRepository) CreateJapanFund(ctx context.Context, dto CreateJapanFundDto) (*model.JapanFund, error) {
	args := m.Called(ctx, dto)
	return args.Get(0).(*model.JapanFund), args.Error(1)
}

func (m *MockJapanFundRepository) DeleteJapanFund(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}