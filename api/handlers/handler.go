package handlers

import (
	"Auth-Service/storage/postgres"

	"go.uber.org/zap"
)

type Handler struct {
	UsersRepo *postgres.UserRepository
	Log          *zap.Logger
}

func NewHandler(users *postgres.UserRepository, log *zap.Logger) *Handler {
	return &Handler{
		UsersRepo: users,
		Log:          log}
}
