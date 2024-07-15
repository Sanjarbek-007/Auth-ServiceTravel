package handlers

import (
	"net/http"

	pb "Auth-Service/genproto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) Register(ctx *gin.Context) {
	var request pb.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	response, err := h.UsersService.Register(ctx, &request)
	if err != nil {
		h.Log.Error("Failed to create user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": response.Id})
}

func (h *Handler) Login(ctx *gin.Context) {
	var request pb.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.Log.Error("Failed to bind JSON", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	response, err := h.UsersService.Login(ctx, &request)
	if err != nil {
		h.Log.Error("Failed to login user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to login user"})
		return
	}

	h.UsersService.Ge
		ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}
