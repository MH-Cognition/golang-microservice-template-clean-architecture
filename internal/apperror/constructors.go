package apperrors

import "errors"

func New(code Code, message string) *AppError {
	if message == "" {
		message = defaultMessage(code)
	}
	return &AppError{Code: code, Message: message}
}

func Wrap(code Code, message string, err error) *AppError {
	if message == "" {
		message = defaultMessage(code)
	}
	return &AppError{Code: code, Message: message, Err: err}
}

func From(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return &AppError{
		Code:    Internal,
		Message: defaultMessage(Internal),
		Err:     err,
	}
}
