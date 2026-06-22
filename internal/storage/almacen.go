package storage

import "proyecto-semestral/internal/models"

type InventarioRepository interface {
	// Inventario
	ListarInventario() []models.Inventario
	BuscarInventarioPorID(id int) (models.Inventario, bool)
	CrearInventario(i models.Inventario) models.Inventario
	ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool)
	BorrarInventario(id int) bool
}

type PublicacionRepository interface {
	// Publicacion
	ListarPublicacion() []models.Publicacion
	BuscarPublicacionPorID(id int) (models.Publicacion, bool)
	CrearPublicacion(p models.Publicacion) models.Publicacion
	ActualizarPublicacion(id int, datos models.Publicacion) (models.Publicacion, bool)
	BorrarPublicacion(id int) bool
}

type UserRepository interface {
	//Usuario
	ListarUsuarios() []models.Usuario
	BuscarUsuarioPorID(id int) (models.Usuario, bool)
	CrearUsuario(u models.Usuario) models.Usuario
	ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool)
	BorrarUsuario(id int) bool
	BuscarUsuarioPorEmail(email string) (models.Usuario, bool)
}

// agregar las demas interfaces de cada entidad

type Almacen interface {
	InventarioRepository

	PublicacionRepository

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
	CrearCalificacion(ca models.Calificacion) models.Calificacion
	ActualizarCalificacion(id int, datos models.Calificacion) (models.Calificacion, bool)
	BorrarCalificacion(id int) bool

	UserRepository

	//Acuerdo
	ListarAcuerdos() []models.Acuerdo
	BuscarAcuerdoPorID(id int) (models.Acuerdo, bool)
	CrearAcuerdo(acuerdo models.Acuerdo) models.Acuerdo
	ActualizarAcuerdo(id int, datos models.Acuerdo) (models.Acuerdo, bool)
	BorrarAcuerdo(id int) bool

	//AcuerdoItem
	ListarAcuerdoItems() []models.AcuerdoItem
	BuscarAcuerdoItemPorID(id int) (models.AcuerdoItem, bool)
	CrearAcuerdoItem(item models.AcuerdoItem) models.AcuerdoItem
	ActualizarAcuerdoItem(id int, datos models.AcuerdoItem) (models.AcuerdoItem, bool)
	BorrarAcuerdoItem(id int) bool
}
