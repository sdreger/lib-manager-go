package errors

import (
	"errors"
	"fmt"
	"github.com/sdreger/lib-manager-go/internal/response"
)

var (
	ErrNotFound = errors.New("the requested resource could not be found")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e ValidationError) ToAPIError() response.APIError {
	return response.APIError{Field: e.Field, Message: e.Message}
}
