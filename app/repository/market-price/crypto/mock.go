package crypto

import (
	"github.com/stretchr/testify/mock"
)

// MockCryptoRepository の定義
type MockCryptoRepository struct {
    mock.Mock
}

// NewMockCryptoRepository は新しい NewMockCryptoRepository を作成し、初期設定を行います。
func NewMockCryptoRepository() *MockCryptoRepository {
	return &MockCryptoRepository{}
}

func (m *MockCryptoRepository) FetchCryptoPrice(code string) (*Crypto, error) {
	args := m.Called(code)
	return args.Get(0).(*Crypto), args.Error(1)
}