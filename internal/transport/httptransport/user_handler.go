package httptransport

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/apperror"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/dto"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/usecase"
)

type UserHandler struct {
	userUC *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUC: uc,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {

	var input dto.CreateUserRequest

	// 1️⃣ Decode request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return apperrors.Wrap(
			apperrors.InvalidJSON,
			"invalid JSON payload",
			err,
		)
	}

	// 2️⃣ Map to usecase input
	ucInput := usecase.CreateUserInput{
		CognitoSub: input.CognitoSub,
		Email:      input.Email,
		Phone:      input.Phone,
		TenantID:   input.TenantID,
		IsActive:   input.IsActive,
	}

	// 3️⃣ Call usecase
	ucOutput, err := h.userUC.CreateUser(r.Context(), ucInput)
	if err != nil {
		// 🔴 DO NOT wrap — preserve original AppError
		return err
	}

	// 4️⃣ Encode response
	resp := dto.UserResponse{
		ID:       ucOutput.ID,
		TenantID: ucOutput.TenantID,
		Email:    ucOutput.Email,
		Phone:    ucOutput.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return apperrors.Wrap(
			apperrors.Internal,
			"failed to encode response",
			err,
		)
	}

	return nil
}

