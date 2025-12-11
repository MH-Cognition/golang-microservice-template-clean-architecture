package domain

import "time"

type User struct {
	ID         string
	CognitoSub string
	Email      string
	Phone      string
	TenantID   string
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
