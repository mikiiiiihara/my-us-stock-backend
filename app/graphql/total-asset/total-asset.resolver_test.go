package totalasset

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTotalAssetService は TotalAssetService のモックです。
type MockTotalAssetService struct {
	mock.Mock
}

func (m *MockTotalAssetService) TotalAssets(ctx context.Context, day int) ([]*generated.TotalAsset, error) {
	args := m.Called(ctx, day)
	return args.Get(0).([]*generated.TotalAsset), args.Error(1)
}

func (m *MockTotalAssetService) UpdateTotalAsset(ctx context.Context, input generated.UpdateTotalAssetInput) (*generated.TotalAsset, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*generated.TotalAsset), args.Error(1)
}

// TotalAssets メソッドのテスト
func TestTotalAssets(t *testing.T) {
	mockService := new(MockTotalAssetService)
	resolver := NewResolver(mockService)

	totalAssets := []*generated.TotalAsset{
		{
			ID:              "1",
			CashJpy:         10000,
			CashUsd:         100,
			Stock:           50000,
			Fund:            30000,
			Crypto:          20000,
			FixedIncomeAsset: 15000,
			CreatedAt:       "2021-01-01",
		},
	}
	mockService.On("TotalAssets", mock.Anything, 30).Return(totalAssets, nil)

	result, err := resolver.TotalAssets(context.Background(), 30)

	assert.NoError(t, err)
	assert.Equal(t, totalAssets, result)

	mockService.AssertExpectations(t)
}

// UpdateTotalAsset メソッドのテスト
func TestUpdateTotalAsset(t *testing.T) {
	mockService := new(MockTotalAssetService)
	resolver := NewResolver(mockService)

	input := generated.UpdateTotalAssetInput{
		ID:      "1",
		CashJpy: 15000,
		CashUsd: 150,
	}
	mockResponse := &generated.TotalAsset{
		ID:              "1",
		CashJpy:         15000,
		CashUsd:         150,
		Stock:           50000,
		Fund:            30000,
		Crypto:          20000,
		FixedIncomeAsset: 15000,
		CreatedAt:       "2021-01-02",
	}
	mockService.On("UpdateTotalAsset", mock.Anything, input).Return(mockResponse, nil)

	result, err := resolver.UpdateTotalAsset(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, mockResponse, result)

	mockService.AssertExpectations(t)
}
