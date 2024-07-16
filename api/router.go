package api

import (
	_ "Auth-Service/api/docs"
	"Auth-Service/api/handlers"
	"Auth-Service/api/middleware"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter @title API Service
// @version 1.0
// @description API service
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(handler *handlers.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// API routes
	user := r.Group("/user")
	{
		user.POST("/register", handler.Register)
		user.POST("/login", handler.Login)
	}
	auth := r.Group("/user")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/refresh-token", handler.RefreshToken)
		auth.POST("/logout", handler.Logout)
	}

	return r
}
