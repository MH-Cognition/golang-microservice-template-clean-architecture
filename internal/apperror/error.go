package apperrors

type AppError struct {
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return string(e.Code) + ": " + e.Message + " | cause=" + e.Err.Error()
	}
	return string(e.Code) + ": " + e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
