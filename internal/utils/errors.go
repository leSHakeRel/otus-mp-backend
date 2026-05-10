package utils

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrInternalServer = errors.New("internal server error")
	ErrValidation     = errors.New("validation error")
	ErrExternalAPI    = errors.New("external api error")
	ErrInvalidInput   = errors.New("invalid input")
	ErrAlreadyExists  = errors.New("already exists")
)

type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code string, message string, details interface{}, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Err:     err,
	}
}

func WrapError(err error, code string, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
