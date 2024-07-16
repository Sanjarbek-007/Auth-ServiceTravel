package main

import (
	"Auth-Service/api"
	"Auth-Service/api/handlers"
	"Auth-Service/genproto"
	"Auth-Service/logger"
	"Auth-Service/service"
	"Auth-Service/storage/postgres"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "Auth-Service/api/docs"
)

// @title Auth Service API
// @version 1.0
// @description This is the API documentation for the Auth Service.
// @host localhost:8081
// @BasePath /
func main() {
	db, err := postgres.ConnectionDb()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	listen, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen on port 50050: %v", err)
	}

	grpcServer := grpc.NewServer()
	userService := &service.UserService{
		UserRepo: &postgres.UserRepository{},
		Log:      &zap.Logger{},
	}
	genproto.RegisterUserServiceServer(grpcServer, userService)

	go func() {
		log.Printf("gRPC server listening on port 50050")
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	handle:=NewHandler()
	router := api.NewRouter(handle)

	log.Print("HTTP server is running on port 8081")
	log.Fatal(router.Run(":8081"))
}

func NewHandler() *handlers.Handler {
	conn, err := grpc.Dial("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("Failed to connect to gRPC server: %v", err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		log.Panicf("Failed to initialize zap logger: %v", err)
	}

	return &handlers.Handler{
		UsersService: genproto.NewUserServiceClient(conn),
		Log:          logger,
	}
}
