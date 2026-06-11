package storage

import "proyecto-semestral/internal/models"

func (a *AlmacenSQLite) ListarAcuerdos() []models.Acuerdo {
	var acuerdos []models.Acuerdo
	a.db.Find(&acuerdos)
	return acuerdos
}

func (a *AlmacenSQLite) BuscarAcuerdoPorID(id int) (models.Acuerdo, bool) {
	var acuerdo models.Acuerdo
	if err := a.db.First(&acuerdo, id).Error; err != nil {
		return models.Acuerdo{}, false
	}
	return acuerdo, true
}

func (a *AlmacenSQLite) CrearAcuerdo(acuerdo models.Acuerdo) models.Acuerdo {
	a.db.Create(&acuerdo) // GORM rellena el ID autogenerado en &p
	return acuerdo
}

func (a *AlmacenSQLite) ActualizarAcuerdo(id int, datos models.Acuerdo) (models.Acuerdo, bool) {
	var existente models.Acuerdo
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Acuerdo{}, false
	}
	a.db.Model(&existente).Updates(datos)
	return existente, true
}

func (a *AlmacenSQLite) BorrarAcuerdo(id int) bool {
	res := a.db.Delete(&models.Acuerdo{}, id)
	return res.RowsAffected > 0
}
