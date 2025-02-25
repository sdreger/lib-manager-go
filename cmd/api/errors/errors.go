package errors

import (
	"errors"
	"fmt"
	"github.com/sdreger/lib-manager-go/internal/response"
	"strings"
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

type ValidationErrors []ValidationError

func (errors ValidationErrors) Error() string {
	builder := strings.Builder{}
	builder.WriteString("validation errors: [")
	for i, err := range errors {
		builder.WriteString(err.Error())
		if i != len(errors)-1 {
			builder.WriteString("; ")
		}
	}
	builder.WriteString("]")
	return builder.String()
}

func (errors ValidationErrors) ToAPIErrors() []response.APIError {
	apiErrors := make([]response.APIError, len(errors))
	for i, err := range errors {
		apiErrors[i] = err.ToAPIError()
	}

	return apiErrors
}
