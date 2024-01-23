package currency

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockCurrencyRepository は CurrencyRepository のモックです。
type MockCurrencyRepository struct {
    mock.Mock
}

func NewMockCurrencyRepository() *MockCurrencyRepository {
	return &MockCurrencyRepository{}
}

// FetchCurrentUsdJpy は CurrencyService のモックメソッドです。
func (m *MockCurrencyRepository) FetchCurrentUsdJpy(ctx context.Context) (float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Error(1)
}