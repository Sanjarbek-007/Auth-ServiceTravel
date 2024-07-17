package service

import (
	pb "Auth-Service/genproto/users"
	"Auth-Service/storage/postgres"
	"context"
)

type UserService struct {
	UserRepo *postgres.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserService(repo *postgres.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (service *UserService) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return service.UserRepo.Register(ctx, in)
}

func (service *UserService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.RegisterResponse, error) {
	return service.UserRepo.Login(ctx, in)
}

func (service *UserService) Profile(ctx context.Context, in *pb.ProfileRequest) (*pb.ProfileResponse, error) {
	return service.UserRepo.Profile(ctx, in)
}

func (service *UserService) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	return service.UserRepo.UpdateProfile(ctx, in)
}

func (service *UserService) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	return service.UserRepo.GetUsers(ctx, in)
}

func (service *UserService) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return service.UserRepo.DeleteUser(ctx, in)
}

func (service *UserService) ResetPassword(ctx context.Context, in *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	return service.UserRepo.ResetPassword(ctx, in)
}

func (service *UserService) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return service.UserRepo.Logout(ctx, in)
}
