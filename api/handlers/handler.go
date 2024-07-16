package handlers

import (
	"Auth-Service/genproto"

	"go.uber.org/zap"
)

type Handler struct {
	UsersService genproto.UserServiceClient
	Log          *zap.Logger
}

