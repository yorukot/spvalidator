package spvalidator

import (
	"errors"
	"fmt"
)

// ValidationError describes a failed validation.
type ValidationError struct {
	Tag     string
	Value   any
	Param   any
	Message string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Param != nil {
		return fmt.Sprintf("%s validation failed for %v with parameter %v", e.Tag, e.Value, e.Param)
	}
	return fmt.Sprintf("%s validation failed for %v", e.Tag, e.Value)
}

// IsValidationError reports whether err contains a ValidationError.
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

func fail(tag string, value any, param any, message string) error {
	return &ValidationError{
		Tag:     tag,
		Value:   value,
		Param:   param,
		Message: message,
	}
}

func failf(tag string, value any, param any, format string, args ...any) error {
	return fail(tag, value, param, fmt.Sprintf(format, args...))
}
