package totalassets

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTotalAssetService の定義
type MockTotalAssetService struct {
    mock.Mock
}

func (m *MockTotalAssetService) CreateTodayTotalAsset(ctx context.Context, c *gin.Context) (string, error) {
    args := m.Called(ctx, c)
    return args.Get(0).(string), args.Error(1)
}

func TestTotalAssetController_CreateTodayTotalAsset(t *testing.T) {
    mockService := new(MockTotalAssetService)
    controller := NewTotalAssetController(mockService)

    expectedMessage := "OK"
    mockService.On("CreateTodayTotalAsset", mock.Anything, mock.AnythingOfType("*gin.Context")).Return(expectedMessage, nil)

    // HTTP リクエストを作成
    req, _ := http.NewRequest(http.MethodPost, "/", nil)

    // gin のコンテキストを作成し、HTTP リクエストを含める
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = req

    // ここで直接メソッドを呼び出す
    controller.CreateTodayTotalAsset(c)

    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), expectedMessage)
    mockService.AssertExpectations(t)
}

