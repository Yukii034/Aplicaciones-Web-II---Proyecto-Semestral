package storage

import "proyecto-semestral/internal/models"

type Almacen2 interface {
	// Reputacion
	ListarReputacion() []models.Reputacion
	BuscarReputacionPorID(id int) (models.Reputacion, bool)
	CrearReputacion(r models.Reputacion) models.Reputacion
	ActualizarReputacion(id int, datos models.Reputacion) (models.Reputacion, bool)
	BorrarReputacion(id int) bool

	// Logro_Usuario
	ListarLogro_Usuario() []models.Logro_Usuario
	BuscarLogro_UsuarioPorID(id int) (models.Logro_Usuario, bool)
	CrearLogro_Usuario(lu models.Logro_Usuario) models.Logro_Usuario
	ActualizarLogro_Usuario(id int, datos models.Logro_Usuario) (models.Logro_Usuario, bool)
	BorrarLogro_Usuario(id int) bool

	//Logro
	ListarLogro() []models.Logro
	BuscarLogroPorID(id int) (models.Logro, bool)
	CrearLogro(l models.Logro) models.Logro
	ActualizarLogro(id int, datos models.Logro) (models.Logro, bool)
	BorrarLogro(id int) bool

	//Calificacion
	ListarCalificacion() []models.Calificacion
	BuscarCalificacionPorID(id int) (models.Calificacion, bool)
	CrearCalificacion(ca models.Calificacion) models.Logro
	ActualizarCalificacion(id int, datos models.Calificacion) (models.Calificacion, bool)
	BorrarCalificacion(id int) bool
}
