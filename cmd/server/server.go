package server

import (
	"Auth-Service/config"
	"Auth-Service/genproto/users"
	"Auth-Service/service"
	"Auth-Service/storage/postgres"
	"log"
	"net"

	"google.golang.org/grpc"
)

func ServerRun(userRepo *postgres.UserRepository, cfg *config.Config) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	users.RegisterUserServiceServer(s, service.NewUserService(userRepo))

	log.Printf("Server is running on %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
