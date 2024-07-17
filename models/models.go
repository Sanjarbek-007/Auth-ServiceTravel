package models

// RegisterRequest represents the registration request payload.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email string    `json:"email" binding:"required"`
	FullName string    `json:"full_name" binding:"required"`
}

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the update profile request payload.
type UpdateProfileRequest struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Bio      string `json:"bio"`
}

// RefreshRequest represents the refresh token request payload.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Tokens represents the tokens response payload.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ProfileResponse represents the profile response payload.
type ProfileResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Bio      string `json:"bio"`
}

type Success struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Failed struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Logout struct{
	Message string `json:"message"`
}