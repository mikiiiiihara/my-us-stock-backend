package user

import (
	"context"
	"my-us-stock-backend/app/graphql/generated"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService の定義
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, name string, email string) (*generated.User, error) {
    args := m.Called(ctx, name, email)
    return args.Get(0).(*generated.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id uint) (*generated.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*generated.User), args.Error(1)
}

// MockUserService が UserService インターフェースを実装することを確認
var _ UserService = (*MockUserService)(nil)

func TestCreateUser(t *testing.T) {
    mockService := new(MockUserService)
    resolver := NewResolver(mockService)

    mockUser := &generated.User{ID: "1", Name: "John Doe", Email: "johndoe@example.com"}
    mockService.On("CreateUser", mock.Anything, "John Doe", "johndoe@example.com").Return(mockUser, nil)

    result, err := resolver.CreateUser(context.Background(), "John Doe", "johndoe@example.com")
    assert.NoError(t, err)
    assert.IsType(t, &generated.User{}, result)
    assert.Equal(t, "John Doe", result.Name)
    assert.Equal(t, "johndoe@example.com", result.Email)
}

func TestUser(t *testing.T) {
    mockService := new(MockUserService)
    resolver := NewResolver(mockService)

    mockUser := &generated.User{ID: "1", Name: "John Doe", Email: "johndoe@example.com"}
    mockService.On("GetUserByID", mock.Anything, uint(1)).Return(mockUser, nil)

    idStr := strconv.FormatUint(uint64(1), 10)
    result, err := resolver.User(context.Background(), idStr)
    assert.NoError(t, err)
    assert.IsType(t, &generated.User{}, result)
    assert.Equal(t, "John Doe", result.Name)
    assert.Equal(t, "johndoe@example.com", result.Email)
}
