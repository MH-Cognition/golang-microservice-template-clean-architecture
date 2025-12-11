package domain

import (
	"errors"
	"time"
)

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

//domain behaviour
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email can't be empty")
	} 

	return nil
}
