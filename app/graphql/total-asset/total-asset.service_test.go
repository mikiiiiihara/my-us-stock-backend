package totalasset

import (
	"context"
	"my-us-stock-backend/app/common/auth"
	"my-us-stock-backend/app/database/model"
	"my-us-stock-backend/app/graphql/generated"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

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
	mockRepo := new(MockTotalAssetRepository)
	mockAuth := auth.NewMockAuthService()
	service := NewTotalAssetService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	mockAssets := []model.TotalAsset{
		{CashJpy: 10000, CashUsd: 100, Stock: 50000},
	}
	mockRepo.On("FetchTotalAssetListById", mock.Anything, userId, 30).Return(mockAssets, nil)

	// テスト対象メソッドの実行
	assets, err := service.TotalAssets(context.Background(), 30)
	assert.NoError(t, err)
	assert.Len(t, assets, 1)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

// UpdateTotalAsset メソッドのテスト
func TestUpdateTotalAssetService(t *testing.T) {
	mockRepo := new(MockTotalAssetRepository)
	mockAuth := auth.NewMockAuthService()
	service := NewTotalAssetService(mockRepo, mockAuth)

	// モックの期待値設定
	userId := uint(1)
	mockAuth.On("FetchUserIdAccessToken", mock.Anything).Return(userId, nil)

	updateInput := generated.UpdateTotalAssetInput{ID: "1", CashJpy: 15000, CashUsd: 150}
	updateDto := totalAssetRepo.UpdateTotalAssetDto{ID: 1, CashJpy: &updateInput.CashJpy, CashUsd: &updateInput.CashUsd}
	mockUpdatedAsset := &model.TotalAsset{CashJpy: 15000, CashUsd: 150}
	mockRepo.On("UpdateTotalAsset", mock.Anything, updateDto).Return(mockUpdatedAsset, nil)

	// テスト対象メソッドの実行
	updatedAsset, err := service.UpdateTotalAsset(context.Background(), updateInput)
	assert.NoError(t, err)
	assert.NotNil(t, updatedAsset)
	assert.Equal(t, float64(15000), updatedAsset.CashJpy)
	assert.Equal(t, float64(150), updatedAsset.CashUsd)

	// モックの呼び出しを検証
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
