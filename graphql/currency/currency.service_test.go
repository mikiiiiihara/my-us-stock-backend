package currency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCurrencyRepository は CurrencyRepository のモックです。
type MockCurrencyRepository struct {
    mock.Mock
}

// FetchCurrentUsdJpy は CurrencyService のモックメソッドです。
func (m *MockCurrencyRepository) FetchCurrentUsdJpy(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}

// TestFetchCurrentUsdJpy は FetchCurrentUsdJpy メソッドのテストです。
func TestFetchCurrentUsdJpy(t *testing.T) {
    // モックの CurrencyService を作成
	mockRepo := new(MockCurrencyRepository)
	service := NewCurrencyService(mockRepo)

    expectedUsdJpy := 133.69
    mockRepo.On("FetchCurrentUsdJpy", mock.Anything).Return(expectedUsdJpy, nil)

    // テストの実行
    result, err := service.FetchCurrentUsdJpy(context.Background())

    // アサーション
    assert.NoError(t, err)
    assert.Equal(t, expectedUsdJpy, result)

    // モックが呼ばれたことを確認
    mockRepo.AssertExpectations(t)
}