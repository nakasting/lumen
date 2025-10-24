package utils

import "fmt"

type ServiceErrorType string

const (
	NotFound ServiceErrorType = "NOT_FOUND"
	Internal ServiceErrorType = "INTERNAL"
	Exists   ServiceErrorType = "EXISTS"
)

type ServiceError struct {
	Type    ServiceErrorType
	Message string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func NewServiceError(t ServiceErrorType, msg string) *ServiceError {
	return &ServiceError{Type: t, Message: msg}
}
