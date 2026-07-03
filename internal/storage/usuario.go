package storage

import (
	"proyecto-semestral/internal/models"

	"gorm.io/gorm"
)

// UsuarioGORM implementa UserRepository sobre GORM.
//
// A diferencia de productos/categorias, los usuarios viven SOLO en GORM: no los
// replicamos en Memoria ni en sqlc (hacerlo seria el mismo ejercicio que ya
// hicimos). El AuthService recibe esta interfaz, nunca este tipo concreto.
type UsuarioGORM struct {
	db *gorm.DB
}

// NuevoUsuarioGORM envuelve una conexion *gorm.DB ya abierta.
func NuevoUsuarioGORM(db *gorm.DB) *UsuarioGORM {
	return &UsuarioGORM{db: db}
}

func (a *UsuarioGORM) ListarUsuarios() []models.Usuario {
	var usuarios []models.Usuario
	a.db.Find(&usuarios)
	return usuarios
}

func (a *UsuarioGORM) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	var usuario models.Usuario
	if err := a.db.First(&usuario, id).Error; err != nil {
		return models.Usuario{}, false
	}
	return usuario, true
}

func (a *UsuarioGORM) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	var u models.Usuario
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil {
		return models.Usuario{}, false
	}
	return u, true
}

func (a *UsuarioGORM) CrearUsuario(usuario models.Usuario) models.Usuario {
	a.db.Create(&usuario) // GORM rellena el ID autogenerado en &p
	return usuario
}

func (a *UsuarioGORM) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	var existente models.Usuario
	if err := a.db.First(&existente, id).Error; err != nil {
		return models.Usuario{}, false
	}
	a.db.Model(&existente).Updates(datos)
	return existente, true
}

func (a *UsuarioGORM) BorrarUsuario(id int) bool {
	res := a.db.Delete(&models.Usuario{}, id)
	return res.RowsAffected > 0
}
