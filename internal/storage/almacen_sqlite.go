package storage

import (
	"gorm.io/gorm"
)

type AlmacenSQLite struct {
	db *gorm.DB
	*UsuarioGORM
}

func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{
		db:          db,
		UsuarioGORM: NuevoUsuarioGORM(db)}
}

// var _ Almacen = (*AlmacenSQLite)(nil)
