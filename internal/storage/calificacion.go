package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarCalificacion() []models.Calificacion {
	var items []models.Calificacion
	a.db.Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarCalificacionPorID(id int) (models.Calificacion, bool) {
	var item models.Calificacion
	if err := a.db.First(&item, id).Error; err != nil {
		return models.Calificacion{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearCalificacion(c models.Calificacion) models.Calificacion {
	a.db.Create(&c) // GORM rellena el ID autogenerado en &p
	return c
}

func (a *AlmacenSQLite) ActualizarCalificacion(id int, datos models.Calificacion) (models.Calificacion, bool) {
	var existente models.Calificacion
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Calificacion{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarCalificacion(id int) bool {
	res := a.db.Delete(&models.Calificacion{}, id)
	return res.RowsAffected > 0
}
