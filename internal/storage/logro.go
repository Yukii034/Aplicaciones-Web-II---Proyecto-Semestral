package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarLogro() []models.Logro {
	var items []models.Logro
	a.db.Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarLogroPorID(id int) (models.Logro, bool) {
	var item models.Logro
	if err := a.db.First(&item, id).Error; err != nil {
		return models.Logro{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearLogro(l models.Logro) models.Logro {
	a.db.Create(&l) // GORM rellena el ID autogenerado en &p
	return l
}

func (a *AlmacenSQLite) ActualizarLogro(id int, datos models.Logro) (models.Logro, bool) {
	var existente models.Logro
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Logro{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarLogro(id int) bool {
	res := a.db.Delete(&models.Logro{}, id)
	return res.RowsAffected > 0
}
