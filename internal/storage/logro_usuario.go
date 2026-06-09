package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarLogro_Usuario() []models.Logro_Usuario {
	var items []models.Logro_Usuario
	a.db.Find(&items)
	return items
}

func (a *AlmacenSQLite) BuscarLogro_UsuarioPorID(id int) (models.Logro_Usuario, bool) {
	var item models.Logro_Usuario
	if err := a.db.First(&item, id).Error; err != nil {
		return models.Logro_Usuario{}, false
	}
	return item, true
}

func (a *AlmacenSQLite) CrearLogro_Usuario(lu models.Logro_Usuario) models.Logro_Usuario {
	a.db.Create(&lu) // GORM rellena el ID autogenerado en &p
	return lu
}

func (a *AlmacenSQLite) ActualizarLogro_Usuario(id int, datos models.Logro_Usuario) (models.Logro_Usuario, bool) {
	var existente models.Logro_Usuario
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Logro_Usuario{}, false
	}
	datos.ID = id
	a.db.Save(&datos)
	return datos, true
}

func (a *AlmacenSQLite) BorrarLogro_Usuario(id int) bool {
	res := a.db.Delete(&models.Logro_Usuario{}, id)
	return res.RowsAffected > 0
}
