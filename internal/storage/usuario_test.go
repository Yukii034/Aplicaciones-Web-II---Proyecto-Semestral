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

// nuevoDBMemoriaUsuario migra el modelo de Usuario a una base de datos SQLite en memoria
func nuevoDBMemoriaUsuario(t *testing.T) *storage.AlmacenSQLite {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migramos el modelo de Usuario
	err = db.AutoMigrate(&models.Usuario{})
	require.NoError(t, err)

	return storage.NuevoAlmacenSQLite(db)
}

func TestAlmacen_CrearYBuscarUsuario(t *testing.T) {
	// Arrange - Configura la DB en memoria
	almacen := nuevoDBMemoriaUsuario(t)

	// Act - Crea un nuevo usuario
	creado := almacen.CrearUsuario(models.Usuario{
		Nombre:       "Juan Pérez",
		Email:        "juan@example.com",
		PasswordHash: "hash_seguro_123",
		Tipo:         "cliente",
		Ciudad:       "Santiago",
		Telefono:     "+56912345678",
		Reputacion:   "buena",
	})

	// Assert - Verifica que se asigne ID y que los datos coincidan al buscar por ID
	assert.NotZero(t, creado.ID, "GORM debe asignar un ID autoincremental")

	encontrado, ok := almacen.BuscarUsuarioPorID(creado.ID)
	require.True(t, ok, "debe encontrar el usuario recién creado por ID")
	assert.Equal(t, "Juan Pérez", encontrado.Nombre)
	assert.Equal(t, "juan@example.com", encontrado.Email)
}

func TestAlmacen_BuscarUsuarioPorEmail(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoriaUsuario(t)
	almacen.CrearUsuario(models.Usuario{
		Nombre: "Maria",
		Email:  "maria@example.com",
		Tipo:   "admin",
	})

	// Act
	encontrado, ok := almacen.BuscarUsuarioPorEmail("maria@example.com")

	// Assert
	require.True(t, ok, "debe encontrar al usuario por su email")
	assert.Equal(t, "Maria", encontrado.Nombre)
}

func TestAlmacen_ListarUsuarios(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoriaUsuario(t)
	almacen.CrearUsuario(models.Usuario{Nombre: "User 1", Email: "u1@test.com", Tipo: "user"})
	almacen.CrearUsuario(models.Usuario{Nombre: "User 2", Email: "u2@test.com", Tipo: "user"})

	// Act
	lista := almacen.ListarUsuarios()

	// Assert
	assert.Len(t, lista, 2, "debe haber exactamente 2 usuarios en la lista")
}

func TestAlmacen_ActualizarUsuario(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoriaUsuario(t)
	creado := almacen.CrearUsuario(models.Usuario{
		Nombre: "Carlos",
		Email:  "carlos@test.com",
		Tipo:   "user",
		Ciudad: "Lima",
	})

	// Act - Modificamos la ciudad y el tipo
	datosNuevos := models.Usuario{
		Ciudad: "Bogotá",
		Tipo:   "premium",
	}
	actualizado, ok := almacen.ActualizarUsuario(creado.ID, datosNuevos)

	// Assert
	require.True(t, ok, "el usuario a actualizar debe existir")
	assert.Equal(t, "Bogotá", actualizado.Ciudad, "la ciudad debió cambiar")
	assert.Equal(t, "premium", actualizado.Tipo, "el tipo debió cambiar")
	assert.Equal(t, "Carlos", actualizado.Nombre, "el nombre debe mantenerse intacto")
}

func TestAlmacen_BorrarUsuario(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoriaUsuario(t)
	creado := almacen.CrearUsuario(models.Usuario{Nombre: "Borrar Me", Email: "borrame@test.com", Tipo: "user"})

	// Act
	eliminado := almacen.BorrarUsuario(creado.ID)

	// Assert
	assert.True(t, eliminado, "debe retornar true indicando que se afectaron filas")

	// Verificar que realmente ya no existe en la DB
	_, ok := almacen.BuscarUsuarioPorID(creado.ID)
	assert.False(t, ok, "el usuario ya no debería existir en el almacén")
}

func TestAlmacen_BuscarUsuarioInexistente(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoriaUsuario(t)

	// Act & Assert para ID inexistente
	_, okID := almacen.BuscarUsuarioPorID(999)
	assert.False(t, okID, "debe retornar false para un ID inexistente")

	// Act & Assert para Email inexistente
	_, okEmail := almacen.BuscarUsuarioPorEmail("no-existe@test.com")
	assert.False(t, okEmail, "debe retornar false para un email inexistente")
}
