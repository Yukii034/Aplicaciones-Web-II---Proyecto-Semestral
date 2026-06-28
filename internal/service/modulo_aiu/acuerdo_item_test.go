package modulo_aiu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	aiu "proyecto-semestral/internal/service/modulo_aiu"
	"proyecto-semestral/internal/storage"
)

// --- Mock Compartido para AcuerdoItem ---
type acuerdoItemRepoMock struct {
	mock.Mock
}

func (m *acuerdoItemRepoMock) CrearAcuerdoItem(a models.AcuerdoItem) models.AcuerdoItem {
	args := m.Called(a)
	return args.Get(0).(models.AcuerdoItem)
}

func (m *acuerdoItemRepoMock) ActualizarAcuerdoItem(id int, datos models.AcuerdoItem) (models.AcuerdoItem, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.AcuerdoItem), args.Bool(1)
}

func (m *acuerdoItemRepoMock) BorrarAcuerdoItem(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

func (m *acuerdoItemRepoMock) ListarAcuerdoItems() []models.AcuerdoItem {
	args := m.Called()
	return args.Get(0).([]models.AcuerdoItem)
}

func (m *acuerdoItemRepoMock) BuscarAcuerdoItemPorID(id int) (models.AcuerdoItem, bool) {
	args := m.Called(id)
	return args.Get(0).(models.AcuerdoItem), args.Bool(1)
}

// Red de seguridad para asegurar que implementa la interfaz del storage de acuerdos
var _ storage.Acuerdo_ItemRepository = (*acuerdoItemRepoMock)(nil)

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
			repo := new(acuerdoItemRepoMock) // Reutiliza el mock definido en acuerdo_test.go

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
	repo := new(acuerdoItemRepoMock)
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
	repo := new(acuerdoItemRepoMock)
	repo.On("BorrarAcuerdoItem", 999).Return(false)

	svc := aiu.NewAcuerdoItemService(repo)

	// Act
	err := svc.BorrarAcuerdoItem(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
