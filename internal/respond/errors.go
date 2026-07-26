package respond

// ErrorResponse defines the standard API's error response format
type ErrorResponse struct {
	Body ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    Code         `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

// Code identifies the application-level error independent of the HTTP status code.
type Code string

type FieldError struct {
	Field string `json:"field"`
	Code  string `json:"code"`
	Param string `json:"param,omitempty"`
}

const (
	CodeValidationError     Code = "validation_error"
	CodeInternalServerError Code = "internal_server_error"
	CodeInvalidJSONError    Code = "invalid_json"
)
