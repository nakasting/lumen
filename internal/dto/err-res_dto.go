package dto

type ErrorResponse struct {
	Message string `json:"errors"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorRes struct {
	Errors []ErrorField `json:"errors"`
}

type ErrorField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func NewErrorRes(fields []ErrorField) *ErrorRes {
	return &ErrorRes{
		Errors: fields,
	}
}
