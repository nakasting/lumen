package dto

type ErrorRes struct {
	Errors []ErrorField `json:"errors"`
}

type ErrorField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
