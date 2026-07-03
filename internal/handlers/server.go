package handlers

import (
	"proyecto-semestral/internal/service"
	aiu "proyecto-semestral/internal/service/modulo_aiu"
	pi "proyecto-semestral/internal/service/modulo_pi" //importar subcarpeta de cada modulo en service
	rlc "proyecto-semestral/internal/service/modulo_rlc"
	"proyecto-semestral/internal/storage"
)

// Server agrupa los servicios de los que dependen los handlers.
//
// Guarda la capa de servicio (no el storage directo): los handlers quedan
// delgados: decodifican el request, llaman al servicio y traducen el resultado
// a HTTP.
type Server struct {
	Almacen       storage.Almacen
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

// Deps agrupa las dependencias requeridas para construir un Server.
//
// Antes NewServer recibia un parametro posicional por servicio; agregar una
// entidad obligaba a cambiar la firma Y todos los call-sites, y dos parametros
// del mismo tipo eran faciles de intercambiar por error. Con un struct de
// dependencias y campos NOMBRADOS, agregar una entidad es agregar un campo:
// nada mas se rompe y desaparece el riesgo de intercambiar argumentos.
type Deps struct {
	Almacen       storage.Almacen
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
}

// NewServer construye un Server a partir de sus dependencias.
func NewServer(d Deps) *Server {
	return &Server{
		Almacen:       d.Almacen,
		Inventario:    d.Inventario,
		Publicacion:   d.Publicacion,
		Reputacion:    d.Reputacion,
		Logro_Usuario: d.Logro_Usuario,
		Logro:         d.Logro,
		Calificacion:  d.Calificacion,
		Acuerdo:       d.Acuerdo,
		AcuerdoItem:   d.AcuerdoItem,
		Usuario:       d.Usuario,
		Auth:          d.Auth,
	}
}
