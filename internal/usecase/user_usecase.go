package usecase

import (
	"context"
	"time"

	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/domain"
	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/repository"
	"github.com/google/uuid"
)

type UserUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(UserRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: UserRepo,
	}
}

type CreateUserInput struct {
	CognitoSub string
	Email      string
	Phone      string
	TenantID   string
	IsActive   bool
}

type CreateUserOutput struct {
	ID         string
	CognitoSub string
	Email      string
	Phone      string
	TenantID   string
}

func (uc *UserUsecase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {

	user := &domain.User{
		ID:         uuid.NewString(),
		CognitoSub: input.CognitoSub,
		Email:      input.Email,
		Phone:      input.Phone,
		TenantID:   input.TenantID,
		IsActive:   input.IsActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdUser, err := uc.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	output := &CreateUserOutput{
		ID:         createdUser.ID,
		CognitoSub: createdUser.CognitoSub,
		Email:      createdUser.Email,
		Phone:      createdUser.Phone,
		TenantID:   createdUser.TenantID,
	}

	return output, nil
}
