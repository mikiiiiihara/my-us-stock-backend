package utils

import "fmt"

type GraphQLAuthError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func (e *GraphQLAuthError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}