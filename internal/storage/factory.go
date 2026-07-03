package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite" // dialector GORM para SQLite (pure-Go)
	"gorm.io/driver/postgres"    // dialector GORM para PostgreSQL
	"gorm.io/gorm"

	"proyecto-semestral/internal/models"
)

// Recursos agrupa todo lo que la capa de almacenamiento expone a la aplicacion:
// el almacen (GORM) y una funcion para cerrar conexiones al apagar.
type Recursos struct {
	Almacen      Almacen
	Usuarios     UserRepository
	BackendUsado string
	Cerrar       func() error
}

// Inicializar centraliza el plumbing de almacenamiento (patron Factory).
//
// El motor de base de datos se elige por configuracion (driver):
//   - "sqlite"   (por defecto): archivo local, ideal para desarrollo.
//   - "postgres": usa el DSN (dsn); es el motor que usa el contenedor Docker.
func Inicializar(driver, dsn, rutaDB string) (*Recursos, error) {
	// 1. GORM es el DUENO DEL ESQUEMA: abre (segun el motor) y migra.
	gdb, err := abrirGorm(driver, dsn, rutaDB)
	if err != nil {
		return nil, err
	}
	if err := gdb.AutoMigrate(
		&models.Inventario{},
		&models.Publicacion{},
		&models.Reputacion{},
		&models.Logro_Usuario{},
		&models.Logro{},
		&models.Calificacion{},
		&models.Usuario{},
		&models.Acuerdo{},
		&models.AcuerdoItem{}); err != nil {
		return nil, fmt.Errorf("AutoMigrate: %w", err)
	}
	almacenGorm := NuevoAlmacenSQLite(gdb)

	// 2. Usuarios viven SIEMPRE en GORM.
	usuarios := NuevoUsuarioGORM(gdb)

	// 3. Cierre ordenado de la conexion GORM.
	cerrar := func() error {
		sqlDB, err := gdb.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return &Recursos{
		Almacen:      almacenGorm,
		Usuarios:     usuarios,
		BackendUsado: "gorm",
		Cerrar:       cerrar,
	}, nil
}

// abrirGorm elige el Dialector segun el driver y abre la conexion.
func abrirGorm(driver, dsn, rutaDB string) (*gorm.DB, error) {
	switch driver {
	case "postgres":
		var gdb *gorm.DB
		var err error
		for intento := 1; intento <= 10; intento++ {
			gdb, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				return gdb, nil
			}
			log.Printf("PostgreSQL no esta listo (intento %d/10): %v", intento, err)
			time.Sleep(2 * time.Second)
		}
		return nil, fmt.Errorf("conectar a PostgreSQL tras reintentos: %w", err)
	default: // "sqlite"
		gdb, err := gorm.Open(sqlite.Open(rutaDB), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("abrir SQLite: %w", err)
		}
		return gdb, nil
	}
}
