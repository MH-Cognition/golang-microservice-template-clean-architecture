package apperrors

type Code string

const (
	InvalidJSON  Code = "INVALID_JSON"
	InvalidInput Code = "INVALID_INPUT"
	NotFound     Code = "NOT_FOUND"
	Conflict     Code = "CONFLICT"
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
	Internal     Code = "INTERNAL_ERROR"

	Timeout      Code = "TIMEOUT"
	Dependency   Code = "DEPENDENCY_ERROR"
)
