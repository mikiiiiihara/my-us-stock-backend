package utils

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

type GraphQLAuthError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e *GraphQLAuthError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// GraphQLの認証エラー
func UnauthenticatedError(message string) *gqlerror.Error {
    return &gqlerror.Error{
        Message: message,
        Extensions: map[string]interface{}{
            "code": "UNAUTHENTICATED",
        },
    }
}
