package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarInventario() []models.Inventario {
	var items []models.Inventario
	a.db.Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarInventarioPorID(id int) (models.Inventario, bool) {
	var item models.Inventario
	if err := a.db.First(&item, id).Error; err != nil {
		return models.Inventario{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearInventario(i models.Inventario) models.Inventario {
	a.db.Create(&i)
	return i
}

func (a *AlmacenSQLite) ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool) {
	var existente models.Inventario
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Inventario{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarInventario(id int) bool {
	res := a.db.Delete(&models.Inventario{}, id)
	return res.RowsAffected > 0
}
