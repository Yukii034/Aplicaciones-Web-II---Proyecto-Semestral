package storage_test

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/storage"
)

func nuevoDBMemoria(t *testing.T) *storage.AlmacenSQLite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&models.Inventario{},
		&models.Publicacion{},
		&models.Logro{},
		&models.Usuario{},
	)
	require.NoError(t, err)

	return storage.NuevoAlmacenSQLite(db)
}

func TestAlmacen_CrearYBuscarInventario(t *testing.T) {
	// Arrange - migra el modelo de inventario a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - agrega un nuevo inventario en el almacen aparte
	creado := almacen.CrearInventario(models.Inventario{
		Nombre:    "Laptop Dell",
		Categoria: "Tecnología",
		Cantidad:  1,
	})

	// Assert - verifica si el nuevo inventario se ha agregado correctamente y con id asignado de GORM
	assert.NotZero(t, creado.ID, "GORM debe asignar un ID")

	encontrado, ok := almacen.BuscarInventarioPorID(creado.ID)
	require.True(t, ok, "debe encontrar el inventario recién creado")
	assert.Equal(t, "Laptop Dell", encontrado.Nombre)
}

func TestAlmacen_ListarInventario(t *testing.T) {
	// Arrange - migra el modelo de inventario a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - crea dos nuevos inventarios de prueba
	almacen.CrearInventario(models.Inventario{Nombre: "Silla", Cantidad: 1})
	almacen.CrearInventario(models.Inventario{Nombre: "Mesa", Cantidad: 2})

	// Assert - verifica si los inventarios creados de test se crearon y se listaron correctamente
	// , ademas de contar si fueron los anteriormente dichos
	lista := almacen.ListarInventario()
	assert.Len(t, lista, 2, "debe haber 2 items en la lista")
}

func TestAlmacen_BuscarInexistente(t *testing.T) {
	// Arrange - migra el modelo de inventario a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - busca en el almacen si existe el id 999
	_, ok := almacen.BuscarInventarioPorID(999)

	// Assert - verifica que haya retornado un false ya que no existe
	assert.False(t, ok, "debe retornar false para un ID inexistente")
}
