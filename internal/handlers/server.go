package handlers

import (
	"proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type Server struct {
	Storage storage.Almacen
	Auth    *service.AuthService
}

func NewServer(s storage.Almacen, auth *service.AuthService) *Server {
	return &Server{Storage: s, Auth: auth}
}
