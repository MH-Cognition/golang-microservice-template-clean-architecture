package apperrors

func defaultMessage(code Code) string {
	switch code {
	case InvalidJSON:
		return "Invalid JSON payload"
	case InvalidInput:
		return "Invalid input"
	case NotFound:
		return "Resource not found"
	case Conflict:
		return "Resource conflict"
	case Unauthorized:
		return "Unauthorized"
	case Forbidden:
		return "Forbidden"
	case Timeout:
		return "Request timed out"
	case Dependency:
		return "Dependency failure"
	default:
		return "Internal server error"
	}
}
