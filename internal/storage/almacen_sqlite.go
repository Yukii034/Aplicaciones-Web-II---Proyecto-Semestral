package storage

import (
	"proyecto-semestral/internal/models"
	"time"

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

	reputaciones := []models.Reputacion{
		{PuntosTotales: 150, Nivel: 2, AcuerdosCompl: 5, CalificacionPromedio: 4.5, UsuarioID: 1},
		{PuntosTotales: 80, Nivel: 1, AcuerdosCompl: 2, CalificacionPromedio: 3, UsuarioID: 2},
		{PuntosTotales: 320, Nivel: 4, AcuerdosCompl: 12, CalificacionPromedio: 4.9, UsuarioID: 3},
	}
	a.db.Create(&reputaciones)

	logro_usuario := []models.Logro_Usuario{
		{FechaDesbl: time.Now(), UsuarioID: 1, LogroID: 1},
		{FechaDesbl: time.Now(), UsuarioID: 2, LogroID: 2},
		{FechaDesbl: time.Now(), UsuarioID: 3, LogroID: 1},
	}
	a.db.Create(&logro_usuario)

	logro := []models.Logro{
		{Nombre: "Primer intercambio", Descripcion: "Completa tu primer trueque", PuntosRequeridos: 50},
		{Nombre: "Donador solidario", Descripcion: "Realizaste tu primera donacion", PuntosRequeridos: 30},
		{Nombre: "Usuario confiable", Descripcion: "Alcanzaste una Calificacion mayor a 4.5", PuntosRequeridos: 100},
	}
	a.db.Create(&logro)

	calificaciones := []models.Calificacion{
		{Comentarios: "Muy buen intercambio, puntual y amable", UsuarioID: 1, AcuerdoID: 1},
		{Comentarios: "Todo perfecto, recomendado", UsuarioID: 2, AcuerdoID: 2},
		{Comentarios: "Buena experiencia en general", UsuarioID: 3, AcuerdoID: 1},
	}
	a.db.Create(&calificaciones)
}

var _ Almacen = (*AlmacenSQLite)(nil)
