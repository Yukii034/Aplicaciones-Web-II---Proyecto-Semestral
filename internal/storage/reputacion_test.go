package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
)

func TestAlmacen_CrearYBuscarReputacion(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoria(t)

	// Act
	creado := almacen.CrearReputacion(models.Reputacion{
		PuntosTotales:        100,
		Nivel:                2,
		AcuerdosCompl:        3,
		CalificacionPromedio: 4.5,
		UsuarioID:            1,
	})

	// Assert
	assert.NotZero(t, creado.ID, "GORM debe asignar un ID")

	encontrado, ok := almacen.BuscarReputacionPorID(creado.ID)
	require.True(t, ok, "debe encontrar la reputación recién creada")
	assert.Equal(t, 100, encontrado.PuntosTotales)
}

func TestAlmacen_ListarReputacion(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoria(t)

	// Act
	almacen.CrearReputacion(models.Reputacion{PuntosTotales: 50, UsuarioID: 1})
	almacen.CrearReputacion(models.Reputacion{PuntosTotales: 80, UsuarioID: 2})

	// Assert
	lista := almacen.ListarReputacion()
	assert.Len(t, lista, 2, "debe haber 2 reputaciones en la lista")
}

func TestAlmacen_BuscarReputacionInexistente(t *testing.T) {
	// Arrange
	almacen := nuevoDBMemoria(t)

	// Act
	_, ok := almacen.BuscarReputacionPorID(999)

	// Assert
	assert.False(t, ok, "debe retornar false para un ID inexistente")
}
