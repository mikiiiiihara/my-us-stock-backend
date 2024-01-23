package currency

import (
	"context"
	"testing"

	"my-us-stock-backend/app/repository/market-price/currency"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestFetchCurrentUsdJpy は FetchCurrentUsdJpy メソッドのテストです。
func TestFetchCurrentUsdJpy(t *testing.T) {
    // モックの CurrencyService を作成
	mockRepo := currency.NewMockCurrencyRepository()
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