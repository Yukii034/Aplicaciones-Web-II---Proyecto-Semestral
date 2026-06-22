package storage

import (
	"proyecto-semestral/internal/models"
)

func (a *AlmacenSQLite) ListarUsuarios() []models.Usuario {
	var usuarios []models.Usuario
	a.db.Find(&usuarios)
	return usuarios
}

func (a *AlmacenSQLite) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	var usuario models.Usuario
	if err := a.db.First(&usuario, id).Error; err != nil {
		return models.Usuario{}, false
	}
	return usuario, true
}

func (a *AlmacenSQLite) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

func (a *AlmacenSQLite) CrearUsuario(usuario models.Usuario) models.Usuario {
	a.db.Create(&usuario) // GORM rellena el ID autogenerado en &p
	return usuario
}

func (a *AlmacenSQLite) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	var existente models.Usuario
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Usuario{}, false
	}
	a.db.Model(&existente).Updates(datos)
	return existente, true
}

func (a *AlmacenSQLite) BorrarUsuario(id int) bool {
	res := a.db.Delete(&models.Usuario{}, id)
	return res.RowsAffected > 0
}
