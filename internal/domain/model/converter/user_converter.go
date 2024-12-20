package converter

import (
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	lastLogin := user.LastLogin.Format(time.RFC3339)
	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Status:    user.Status,
		LastLogin: &lastLogin,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func LoginToTokenResponse(accessToken, refreshToken string) *model.TokenResponse {
	return &model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
