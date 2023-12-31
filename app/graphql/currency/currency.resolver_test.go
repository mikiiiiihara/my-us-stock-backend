package currency

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCurrencyService は CurrencyService のモックです。
type MockCurrencyService struct {
    mock.Mock
}

// FetchCurrentUsdJpy は CurrencyService のモックメソッドです。
func (m *MockCurrencyService) FetchCurrentUsdJpy(ctx context.Context) (float64, error) {
    args := m.Called(ctx)
    return args.Get(0).(float64), args.Error(1)
}

// TestGetCurrentUsdJpy は GetCurrentUsdJpy メソッドのテストです。
func TestGetCurrentUsdJpy(t *testing.T) {
    // モックの CurrencyService を作成
    mockService := new(MockCurrencyService)
    resolver := NewResolver(mockService)

    expectedUsdJpy := 133.69
    mockService.On("FetchCurrentUsdJpy", mock.Anything).Return(expectedUsdJpy, nil)

    // テストの実行
    result, err := resolver.CurrentUsdJpy(context.Background())

    // アサーション
    assert.NoError(t, err)
    assert.Equal(t, expectedUsdJpy, result)

    // モックが呼ばれたことを確認
    mockService.AssertExpectations(t)
}
