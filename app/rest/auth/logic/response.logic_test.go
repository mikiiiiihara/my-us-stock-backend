package logic

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendResponse(t *testing.T) {
    rl := NewResponseLogic()
    w := httptest.NewRecorder()
    response := []byte(`{"message":"success"}`)
    code := http.StatusOK

    rl.SendResponse(w, response, code)

    result := w.Result()
    defer result.Body.Close()

    assert.Equal(t, code, result.StatusCode)
    assert.Equal(t, "application/json", result.Header.Get("Content-type"))
    // レスポンスボディの検証が必要な場合は追加する
}

func TestSendNotBodyResponse(t *testing.T) {
    rl := NewResponseLogic()
    w := httptest.NewRecorder()

    rl.SendNotBodyResponse(w)

    result := w.Result()
    defer result.Body.Close()

    assert.Equal(t, http.StatusNoContent, result.StatusCode)
    assert.Equal(t, "application/json", result.Header.Get("Content-type"))
}

func TestCreateErrorResponse(t *testing.T) {
    rl := NewResponseLogic()
    err := errors.New("test error")
    expected := `{"error":"test error"}`

    response := rl.CreateErrorResponse(err)

    assert.JSONEq(t, expected, string(response))
}

func TestCreateErrorStringResponse(t *testing.T) {
    rl := NewResponseLogic()
    errMessage := "test error message"
    expected := `{"error":"test error message"}`

    response := rl.CreateErrorStringResponse(errMessage)

    assert.JSONEq(t, expected, string(response))
}
