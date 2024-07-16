package api

import (
	"Auth-Service/api/handlers"
	"Auth-Service/api/middleware"
	a "Auth-Service/genproto"
  _ "Auth-Service/api/docs"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

// NewRouter @title API Service
// @version 1.0
// @description API service
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(conn *grpc.ClientConn) *gin.Engine {
  r := gin.Default()
  r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

  // Initialize gRPC client for UserService
  userService := a.NewUserServiceClient(conn)
  handler := handlers.Handler{UsersService: userService}

  // API routes
  auth := r.Group("/")
  auth.Use(middleware.AuthMiddleware())
  {
    auth.POST("/register", handler.Register)
    auth.GET("/profile", handler.Profile)
    auth.PUT("/update-profile", handler.UpdateProfile)
    auth.POST("/refresh-token", handler.RefreshToken)
    auth.POST("/logout", handler.Logout)
  }

  return r
}