package apperrors

import "net/http"

func HTTPStatus(code Code) int {
	switch code {
	case InvalidJSON, InvalidInput:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	case Timeout:
		return http.StatusGatewayTimeout
	case Dependency:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
