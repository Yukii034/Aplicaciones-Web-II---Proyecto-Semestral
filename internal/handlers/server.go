package handlers

import (
	"proyecto-semestral/internal/service"
	pi "proyecto-semestral/internal/service/modulo_pi" //importar subcarpeta de cada modulo en service
)

type Server struct {
	Inventario  *pi.InventarioService
	Publicacion *pi.PublicacionService
	// demás servicios de entidades
	Auth *service.AuthService
}

func NewServer(inv *pi.InventarioService, pub *pi.PublicacionService, auth *service.AuthService) *Server {
	return &Server{
		Inventario:  inv,
		Publicacion: pub,
		// añadir cada nuevo servicio
		Auth: auth,
	}
}
