package httptransport

import (
	"encoding/json"
	"net/http"

	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/transport/dto"
	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/usecase"
)

type UserHandler struct {
	userUC *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUC: uc,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {

	var input dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid Input"+err.Error(), http.StatusBadRequest)
		return
	}

	ucInput := &usecase.CreateUserInput{
		CognitoSub: input.CognitoSub,
		Email:      input.Email,
		Phone:      input.Phone,
		TenantID:   input.TenantID,
		IsActive:   input.IsActive,
	}

	ucOutput, err := h.userUC.CreateUser(r.Context(), *ucInput)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	resp := &dto.UserResponse{
		ID:       ucOutput.ID,
		TenantID: ucOutput.TenantID,
		Email:    ucOutput.Email,
		Phone:    ucOutput.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)

}
