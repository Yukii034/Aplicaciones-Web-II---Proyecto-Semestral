package modulo_aiu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	aiu "proyecto-semestral/internal/service/modulo_aiu"
)

// --- Tests para AcuerdoItem ---

// Test con estructura Table Driven para la creación de un AcuerdoItem
func TestAcuerdoService_CrearItem(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.AcuerdoItem
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre: "AcuerdoID inválido (0) -> ErrVacio",
			entrada: models.AcuerdoItem{
				AcuerdoID:    0,
				InventarioID: 45,
				Rol:          "Ofertado",
			},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre: "InventarioID inválido (0) -> ErrVacio",
			entrada: models.AcuerdoItem{
				AcuerdoID:    12,
				InventarioID: 0,
				Rol:          "Solicitado",
			},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre: "Rol vacío -> ErrVacio",
			entrada: models.AcuerdoItem{
				AcuerdoID:    12,
				InventarioID: 45,
				Rol:          "",
			},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre: "AcuerdoItem válido -> sin error y se persiste",
			entrada: models.AcuerdoItem{
				AcuerdoID:    12,
				InventarioID: 45,
				Rol:          "Ofertado",
			},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			repo := new(acuerdoRepoMock) // Reutiliza el mock definido en acuerdo_test.go

			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearAcuerdoItem", c.entrada).Return(guardado)
			}

			svc := aiu.NewAcuerdoItemService(repo)

			// Act
			creado, err := svc.CrearAcuerdoItem(c.entrada)

			// Assert
			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearAcuerdoItem")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearAcuerdoItem", c.entrada)
			}
		})
	}
}

// Test para buscar un AcuerdoItem que no existe
func TestAcuerdoService_BuscarItem_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(acuerdoRepoMock)
	repo.On("BuscarAcuerdoItemPorID", 999).Return(models.AcuerdoItem{}, false)

	svc := aiu.NewAcuerdoItemService(repo)

	// Act
	_, err := svc.BuscarAcuerdoItem(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

// Test para borrar un AcuerdoItem que no existe
func TestAcuerdoService_BorrarItem_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(acuerdoRepoMock)
	repo.On("BorrarAcuerdoItem", 999).Return(false)

	svc := aiu.NewAcuerdoItemService(repo)

	// Act
	err := svc.BorrarAcuerdoItem(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
