package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
)

func TestAlmacen_CrearYBuscarPublicacion(t *testing.T) {
	// Arrange - migra el modelo de publicacion a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - agrega una nueva publicacion en el almacen aparte
	creado := almacen.CrearPublicacion(models.Publicacion{
		Titulo:            "Cambio laptop por tablet",
		TipoOferta:        "intercambio",
		EstadoPublicacion: "disponible",
	})

	// Assert - verifica si la nueva publicacion se ha agregado correctamente y con id asignado de GORM
	assert.NotZero(t, creado.ID, "GORM debe asignar un ID")

	encontrado, ok := almacen.BuscarPublicacionPorID(creado.ID)
	require.True(t, ok, "debe encontrar la publicacion recién creada")
	assert.Equal(t, "Cambio laptop por tablet", encontrado.Titulo)
}

func TestAlmacen_ListarPublicacion(t *testing.T) {
	// Arrange - migra el modelo de publicacion a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - crea dos nuevos inventarios de prueba
	almacen.CrearPublicacion(models.Publicacion{Titulo: "Dono microondas", TipoOferta: "donacion"})
	almacen.CrearPublicacion(models.Publicacion{Titulo: "Cambio bicicleta", TipoOferta: "intercambio"})

	// Assert - verifica si las publicaciones creados de test se crearon y se listaron correctamente
	// , ademas de contar si fueron los anteriormente dichos
	lista := almacen.ListarPublicacion()
	assert.Len(t, lista, 2, "debe haber 2 items en la lista")
}

func TestAlmacen_BuscarInexistenteP(t *testing.T) {
	// Arrange - migra el modelo de publicacion a una memoria aparte
	almacen := nuevoDBMemoria(t)

	// Act - busca en el almacen si existe el id 999
	_, ok := almacen.BuscarPublicacionPorID(999)

	// Assert - verifica que haya retornado un false ya que no existe
	assert.False(t, ok, "debe retornar false para un ID inexistente")
}
