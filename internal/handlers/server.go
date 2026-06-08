package handlers

import "proyecto-semestral/internal/storage"

type Server struct {
	Storage storage.Almacen
}

func NewServer(s storage.Almacen) *Server {
	return &Server{Storage: s}
}
