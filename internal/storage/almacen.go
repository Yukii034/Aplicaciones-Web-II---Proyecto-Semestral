package storage

import "proyecto-semestral/internal/models"

type Almacen interface {
	// Inventario
	ListarInventario() []models.Inventario
	BuscarInventarioPorID(id int) (models.Inventario, bool)
	CrearInventario(i models.Inventario) models.Inventario
	ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool)
	BorrarInventario(id int) bool

	// Publicacion
	ListarPublicacion() []models.Publicacion
	BuscarPublicacionPorID(id int) (models.Publicacion, bool)
	CrearPublicacion(p models.Publicacion) models.Publicacion
	ActualizarPublicacion(id int, datos models.Publicacion) (models.Publicacion, bool)
	BorrarPublicacion(id int) bool
}
