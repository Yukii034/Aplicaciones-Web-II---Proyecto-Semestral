package handlers

import (
	"proyecto-semestral/internal/service"
	pi "proyecto-semestral/internal/service/modulo_pi" //importar subcarpeta de cada modulo en service
	rlc "proyecto-semestral/internal/service/modulo_rlc"
)

type Server struct {
	Inventario    *pi.InventarioService
	Publicacion   *pi.PublicacionService
	Reputacion    *rlc.ReputacionService
	Logro_Usuario *rlc.Logro_UsuarioService
	Logro         *rlc.LogroService
	Calificacion  *rlc.CalificacionService
	Auth          *service.AuthService
	// demás servicios de entidades
}

func NewServer(
	inv *pi.InventarioService,
	pub *pi.PublicacionService,
	rep *rlc.ReputacionService,
	lu *rlc.Logro_UsuarioService,
	l *rlc.LogroService,
	ca *rlc.CalificacionService,
	auth *service.AuthService) *Server {
	return &Server{
		Inventario:    inv,
		Publicacion:   pub,
		Reputacion:    rep,
		Logro_Usuario: lu,
		Logro:         l,
		Calificacion:  ca,
		Auth:          auth,
		// añadir cada nuevo servicio
	}
}
