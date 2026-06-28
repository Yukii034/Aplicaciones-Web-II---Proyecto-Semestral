package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
)

func TestAlmacen_CrearYBuscarLogro(t *testing.T) {
	// Arrange - migra el modelo de logro a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - agrega un nuevo logro en el almacen aparte
	creado := almacen.CrearLogro(models.Logro{
		Nombre:           "Primer Intercambio",
		Descripcion:      "Completaste tu primer trueque",
		PuntosRequeridos: 50,
	})

	// Assert - verifica si el nuevo logro se ha agregado correctamente y con id asignado de GORM
	assert.NotZero(t, creado.ID, "GORM debe asignar un ID")

	encontrado, ok := almacen.BuscarLogroPorID(creado.ID)
	require.True(t, ok, "debe encontrar el logro recién creado")
	assert.Equal(t, "Primer Intercambio", encontrado.Nombre)
}

func TestAlmacen_ListarLogro(t *testing.T) {
	// Arrange - migra el modelo de logro a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - crea dos nuevos logros de prueba
	almacen.CrearLogro(models.Logro{Nombre: "Logro A", PuntosRequeridos: 10})
	almacen.CrearLogro(models.Logro{Nombre: "Logro B", PuntosRequeridos: 20})

	// Assert - verifica si los logros creados de test se crearon y se listaron correctamente
	lista := almacen.ListarLogro()
	assert.Len(t, lista, 2, "debe haber 2 items en la lista")
}

func TestAlmacen_BuscarLogroInexistente(t *testing.T) {
	// Arrange - migra el modelo de logro a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - busca en el almacen si existe el id 999
	_, ok := almacen.BuscarLogroPorID(999)

	// Assert - verifica que haya retornado un false ya que no existe
	assert.False(t, ok, "debe retornar false para un ID inexistente")
}
