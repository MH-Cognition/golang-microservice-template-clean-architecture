package repository

import (
	"context"

	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}
