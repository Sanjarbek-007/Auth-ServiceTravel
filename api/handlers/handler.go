package handlers

import (
	"Auth-Service/genproto"

	"go.uber.org/zap"
)

type Handler struct {
	UsersService genproto.UserServiceClient
	Log          *zap.Logger
}

func NewHandler(user genproto.UserServiceClient, l *zap.Logger) *Handler {
	return &Handler{
		UsersService: user,
		Log:          l,
	}
}
