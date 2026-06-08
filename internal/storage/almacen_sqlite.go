package storage

import "gorm.io/gorm"

type AlmacenSQLite struct {
	db *gorm.DB
}

func NuevoAlmacenSQLite(db *gorm.DB) *AlmacenSQLite {
	return &AlmacenSQLite{db: db}
}

var _ Almacen = (*AlmacenSQLite)(nil)
