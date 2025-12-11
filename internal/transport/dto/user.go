package dto

type CreateUserRequest struct {
	CognitoSub string `json:"cognito_sub"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	TenantID   string `json:"tenant_id"`
	IsActive   bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	CognitoSub *string `json:"cognito_sub,omitempty"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
}

type UserResponse struct {
	ID         string `json:"id"`
	CognitoSub string `json:"cognito_sub"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	TenantID   string `json:"tenant_id"`
	IsActive   bool   `json:"is_active"`
}
