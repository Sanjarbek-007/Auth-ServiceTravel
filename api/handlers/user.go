package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"Auth-Service/genproto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	t "Auth-Service/api/token"
	"go.uber.org/zap"
	"Auth-Service/models"
)
// Register handles user registration.
// @Summary Register a new user
// @Description Register a new user with username and password and email
// @Accept json
// @Produce json
// @Param input body genproto.RegisterRequest true "Registration details" 
// @Success 201 {object} genproto.RegisterResponse
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal status error"
// @Router /user/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var request models.RegisterRequest
	if h.Log == nil {
        fmt.Println("Logger is nil")
        ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Internal server error"})
        return
    }

    if h.UsersService == nil {
        h.Log.Error("UsersService is nil")
        ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Internal server error"})
        return
    }
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Failed{Message: "Invalid request payload", Error: err.Error()})
		return
	}

	response, err := h.UsersService.Register(ctx, &genproto.RegisterRequest{
		Username: request.Username,
		Password: request.Password,
		Email: request.Email,
		FullName: request.FullName,
	})
	if err != nil {
		h.Log.Error("Failed to create user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Failed to create user", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, models.Success{Message: "User created successfully", Data: map[string]string{"user_id": response.Id}})
}


// Login handles user login.
// @Summary Login a user
// @Description Login a user with username and password
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "Login details"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} models.Failed
// @Failure 500 {object} models.Failed
// @Router /user/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var request models.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Failed{Message: "Invalid request payload", Error: err.Error()})
		return
	}

	response, err := h.UsersService.Login(ctx, &genproto.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		h.Log.Error("Failed to login user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Failed to login user", Error: err.Error()})
		return
	}

	tokens, err := t.GENERATEJWTToken(response)
	if err != nil {
		h.Log.Error("Failed to generate token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Failed to generate token", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Tokens{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken})
}

// RefreshToken refreshes user token with refresh token.
// @Summary Refresh user token
// @Description Refresh user token with refresh token
// @Accept json
// @Produce json
// @Param input body models.RefreshRequest true "Refresh token details"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} models.Failed
// @Failure 500 {object} models.Failed
// @Router /refresh-token [post]
func (h *Handler) RefreshToken(ctx *gin.Context) {
	var request models.RefreshRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Failed{Message: "Invalid request payload", Error: err.Error()})
		return
	}

	response, err := h.UsersService.Refresh(ctx, &genproto.RefreshRequest{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		h.Log.Error("Failed to refresh token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Failed to refresh token", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Logout logs out a user by invalidating token.
// @Summary Logout a user
// @Description Logout a user by invalidating token
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Logout
// @Failure 401 {object} models.Failed
// @Router /logout [post]
func (h *Handler) Logout(ctx *gin.Context) {
	_, err := h.ValidateToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Failed{Message: "Unauthorized", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ValidateToken validates the JWT token from Authorization header.
func (h *Handler) ValidateToken(ctx *gin.Context) (jwt.MapClaims, error) {
	tokenString := ctx.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("salom"), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
