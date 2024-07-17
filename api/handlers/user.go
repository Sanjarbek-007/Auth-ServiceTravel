package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"Auth-Service/api/token"
	"Auth-Service/genproto/users"
	"Auth-Service/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Register handles user registration.
// @Summary Register a new user
// @Description Register a new user with username and password and email
// @Security ApiKeyAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body users.RegisterRequest true "Registration details"
// @Success 201 {object} users.RegisterResponse
// @Failure 400 {object} string "bad request"
// @Failure 500 {object} string "internal status error"
// @Router /auth/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	var request models.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Failed{Message: "Invalid request payload", Error: err.Error()})
		return
	}

	response, err := h.UsersRepo.Register(ctx, &users.RegisterRequest{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
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
// @Security ApiKeyAuth
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.LoginRequest true "Login details"
// @Success 200 {object} models.Tokens
// @Failure 400 {object} models.Failed
// @Failure 500 {object} models.Failed
// @Router /auth/login [post]
func (h Handler) Login(ctx *gin.Context) {
	h.Log.Info("Login is working")
	req := users.LoginRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
	}

	res, err := h.UsersRepo.Login(ctx, &req)
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(500, gin.H{"error2": err.Error()})
		return
	}
	var toke users.Token
	err = token.GeneratedAccessJWTToken(res, &toke)

	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(500, gin.H{"error3": err.Error()})
	}
	err = token.GeneratedRefreshJWTToken(res, &toke)
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(500, gin.H{"error4": err.Error()})
	}

	ctx.JSON(http.StatusOK, &toke)
	h.Log.Info("login is succesfully ended")

}

// @Summary Refresh token
// @Description it changes your access token
// @Security ApiKeyAuth
// @Tags auth
// @Param userinfo body users.CheckRefreshTokenRequest true "token"
// @Success 200 {object} users.Token
// @Failure 400 {object} string "Invalid date"
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /user/refresh [post]
func (h Handler) Refresh(ctx *gin.Context) {
	h.Log.Info("Refresh is working")
	req := users.CheckRefreshTokenRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	_, err := token.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id, err := token.GetUserIdFromRefreshToken(req.RefreshToken)
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	res := users.Token{RefreshToken: req.RefreshToken}

	err = token.GeneratedAccessJWTToken(&users.RegisterResponse{Id: id}, &res)
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	ctx.JSON(http.StatusOK, &res)
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

// Profile retrieves user profile details.
// @Summary Get user profile
// @Description Retrieve user profile details
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.ProfileResponse
// @Failure 401 {object} models.Failed
// @Failure 500 {object} models.Failed
// @Router /user/profile/{user_id} [get]
func (h *Handler) Profile(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	fmt.Println(userID)

	claims, err := h.ValidateToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Failed{Message: "Unauthorized", Error: err.Error()})
		return
	}

	if claims["sub"] != userID {
		ctx.JSON(http.StatusUnauthorized, models.Failed{Message: "Unauthorized access to profile"})
		return
	}

	request := &users.ProfileRequest{UserId: userID}
	response, err := h.UsersRepo.Profile(ctx, request)
	if err != nil {
		h.Log.Error("Failed to fetch profile", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Failed to fetch profile", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}


// UpdateProfile updates user profile details.
// @Summary Update user profile
// @Description Update user profile details
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param input body models.UpdateProfileRequest true "Update details"
// @Success 200 {object} models.ProfileResponse
// @Failure 400 {object} models.Failed
// @Failure 401 {object} models.Failed
// @Failure 500 {object} models.Failed
// @Router /user/profileUpdate/{user_id} [put]
func (h *Handler) UpdateProfile(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	claims, err := h.ValidateToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Failed{Message: "Unauthorized", Error: err.Error()})
		return
	}

	if claims["sub"] != userID {
		ctx.JSON(http.StatusUnauthorized, models.Failed{Message: "Unauthorized access to profile"})
		return
	}

	var request models.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, models.Failed{Message: "Invalid request payload", Error: err.Error()})
		return
	}

	request.UserID = userID
	response, err := h.UsersRepo.UpdateProfile(ctx, &users.UpdateProfileRequest{
		Id:       request.UserID,
		FullName: request.FullName,
		Bio:      request.Bio,
	})
	if err != nil {
		h.Log.Error("Failed to update profile", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, models.Failed{Message: "Failed to update profile", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary delete user
// @Description you can delete your profile
// @Security ApiKeyAuth
// @Tags User
// @Param user_id path string true "user_id"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /user/users/{user_id} [delete]
func (h Handler) Delete(ctx *gin.Context) {
	h.Log.Info("Delete is working")
	id := ctx.Param("user_id")
	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user id is incorrect"})
		return
	}

	_, err = h.UsersRepo.DeleteUser(ctx, &users.DeleteUserRequest{Id: id})
	if err != nil {
		h.Log.Error(err.Error())
		ctx.JSON(500, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	h.Log.Info("Delete ended")
}
