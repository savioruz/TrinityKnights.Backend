package model

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,min=8,lte=255"`
	Name     string `json:"name" validate:"required,lte=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,min=8,lte=255"`
}

// UpdateRequest Just update profile
type UpdateRequest struct {
	Name string `json:"name" validate:"omitempty,lte=100"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,min=8,lte=255"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
