package api

import (
	_ "Auth-Service/api/docs"
	"Auth-Service/api/handlers"
	"Auth-Service/api/middleware"

	// "Auth-Service/api/middleware"

	// "Auth-Service/api/middleware"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter @title API Service
// @version 1.0
// @description API service
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(handler *handlers.Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// API routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)
	}
	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/profile/:user_id", handler.Profile)
		user.PUT("/profileUpdate/:user_id", handler.UpdateProfile)
		user.DELETE("/users/:user_id", handler.Delete)
		user.POST("/user/:user_id/follow", handler.FollowUser)
		user.GET("/user/:user_id/followers", handler.FollowersUsers)
	}

	return r
}
