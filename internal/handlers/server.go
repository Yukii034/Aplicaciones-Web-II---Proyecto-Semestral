package handlers

import (
	"proyecto-semestral/internal/service"
	aiu "proyecto-semestral/internal/service/modulo_aiu"
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
	Acuerdo       *aiu.AcuerdoService
	AcuerdoItem   *aiu.AcuerdoItemService
	Usuario       *aiu.UsuarioService
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
	ac *aiu.AcuerdoService,
	acI *aiu.AcuerdoItemService,
	user *aiu.UsuarioService,
	auth *service.AuthService) *Server {

	return &Server{
		Inventario:    inv,
		Publicacion:   pub,
		Reputacion:    rep,
		Logro_Usuario: lu,
		Logro:         l,
		Calificacion:  ca,
		Acuerdo:       ac,
		AcuerdoItem:   acI,
		Usuario:       user,
		Auth:          auth,
		// añadir cada nuevo servicio
	}
}
