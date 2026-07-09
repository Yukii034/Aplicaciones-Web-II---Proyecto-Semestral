package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarPublicacion() []models.Publicacion {
	var items []models.Publicacion
	a.db.Preload("Inventario").Preload("Usuario").Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarPublicacionPorID(id int) (models.Publicacion, bool) {
	var item models.Publicacion
	if err := a.db.Preload("Inventario").Preload("Usuario").First(&item, id).Error; err != nil {
		return models.Publicacion{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearPublicacion(p models.Publicacion) models.Publicacion {
	a.db.Create(&p)
	// recarga con las relaciones después de crear
	a.db.Preload("Inventario").Preload("Usuario").First(&p, p.ID)
	return p
}

func (a *AlmacenSQLite) ActualizarPublicacion(id int, datos models.Publicacion) (models.Publicacion, bool) {
	var existente models.Publicacion
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Publicacion{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	// recarga con relaciones
	a.db.Preload("Inventario").Preload("Usuario").First(&datos, id)
	return datos, true
}

func (a *AlmacenSQLite) BorrarPublicacion(id int) bool {
	res := a.db.Delete(&models.Publicacion{}, id)
	return res.RowsAffected > 0
}
