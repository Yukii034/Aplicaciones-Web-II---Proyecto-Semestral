package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarAcuerdoItems() []models.AcuerdoItem {
	var items []models.AcuerdoItem
	a.db.Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarAcuerdoItemPorID(id int) (models.AcuerdoItem, bool) {
	var item models.AcuerdoItem
	if err := a.db.First(&item, id).Error; err != nil {
		return models.AcuerdoItem{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearAcuerdoItem(item models.AcuerdoItem) models.AcuerdoItem {
	a.db.Create(&item) // GORM rellena el ID autogenerado en &p
	return item
}

func (a *AlmacenSQLite) ActualizarAcuerdoItem(id int, datos models.AcuerdoItem) (models.AcuerdoItem, bool) {
	var existente models.AcuerdoItem
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.AcuerdoItem{}, false
	}
	a.db.Model(&existente).Updates(datos)
	return existente, true
}

func (a *AlmacenSQLite) BorrarAcuerdoItem(id int) bool {
	res := a.db.Delete(&models.AcuerdoItem{}, id)
	return res.RowsAffected > 0
}
