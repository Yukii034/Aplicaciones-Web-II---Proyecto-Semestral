package storage

import (
	"proyecto-semestral/internal/models"

	"gorm.io/gorm"
)

type AlmacenSQLite struct {
	db *gorm.DB
}

func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

func (a *AlmacenSQLite) SembrarSiVacio() {
	var n int64
	a.db.Model(&models.Inventario{}).Count(&n)
	if n > 0 {
		return
	}

	inventarios := []models.Inventario{
		{Nombre: "Bicicleta de montaña", Descripcion: "Poco uso, en buen estado", Categoria: "Deportes", EstadoObjeto: "usado", Disponibilidad: "disponible", Cantidad: 1, FotoUrl: ""},
		{Nombre: "Microondas Samsung", Descripcion: "Funciona perfecto", Categoria: "Electrodomésticos", EstadoObjeto: "usado", Disponibilidad: "disponible", Cantidad: 1, FotoUrl: ""},
		{Nombre: "Libros universitarios", Descripcion: "Cálculo y física", Categoria: "Educación", EstadoObjeto: "nuevo", Disponibilidad: "disponible", Cantidad: 5, FotoUrl: ""},
	}
	a.db.Create(&inventarios)

	publicaciones := []models.Publicacion{
		{Titulo: "Cambio bicicleta por tablet", TipoOferta: "intercambio", EstadoPublicacion: "disponible", Mensaje: "Interesados escribir", UsuarioID: 1, InventarioID: 1},
		{Titulo: "Dono microondas", TipoOferta: "donacion", EstadoPublicacion: "disponible", Mensaje: "Pasar a retirar", UsuarioID: 1, InventarioID: 2},
	}
	a.db.Create(&publicaciones)
}

var _ Almacen = (*AlmacenSQLite)(nil)
