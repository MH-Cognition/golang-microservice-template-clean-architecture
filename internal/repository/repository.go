package repository

import (
	"context"

	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
