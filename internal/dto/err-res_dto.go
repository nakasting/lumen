package dto

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{message}
}

type ValidationErrors struct {
	Errors []FieldError `json:"errors"`
}

func NewValidationErrors(fieldErrors []FieldError) *ValidationErrors {
	return &ValidationErrors{Errors: fieldErrors}
}

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func NewFieldError(field string, error string) *FieldError {
	return &FieldError{field, error}
}
