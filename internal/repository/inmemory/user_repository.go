package inmemory

import (
	"context"
	"sync"

	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/domain"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users []*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: []*domain.User{},
	}
}

func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users = append(r.users, user)
	return user, nil
}
