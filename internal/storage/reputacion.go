package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarReputacion() []models.Reputacion {
	var items []models.Reputacion
	a.db.Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarReputacionPorID(id int) (models.Reputacion, bool) {
	var item models.Reputacion
	if err := a.db.First(&item, id).Error; err != nil {
		return models.Reputacion{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearReputacion(r models.Reputacion) models.Reputacion {
	a.db.Create(&r) // GORM rellena el ID autogenerado en &p
	return r
}

func (a *AlmacenSQLite) ActualizarReputacion(id int, datos models.Reputacion) (models.Reputacion, bool) {
	var existente models.Reputacion
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Reputacion{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarReputacion(id int) bool {
	res := a.db.Delete(&models.Reputacion{}, id)
	return res.RowsAffected > 0
}
