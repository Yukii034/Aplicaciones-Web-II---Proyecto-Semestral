package modulo_pi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	pi "proyecto-semestral/internal/service/modulo_pi"
	"proyecto-semestral/internal/storage"
)

// --- Mock ---
type inventarioRepoMock struct {
	mock.Mock
}

func (m *inventarioRepoMock) ListarInventario() []models.Inventario {
	args := m.Called()
	return args.Get(0).([]models.Inventario)
}

func (m *inventarioRepoMock) BuscarInventarioPorID(id int) (models.Inventario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Inventario), args.Bool(1)
}

func (m *inventarioRepoMock) CrearInventario(i models.Inventario) models.Inventario {
	args := m.Called(i)
	return args.Get(0).(models.Inventario)
}

func (m *inventarioRepoMock) ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Inventario), args.Bool(1)
}

func (m *inventarioRepoMock) BorrarInventario(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// red de seguridad
var _ storage.InventarioRepository = (*inventarioRepoMock)(nil)

// --- Tests ---
func TestInventarioService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Inventario
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "nombre vacio -> ErrVacio",
			entrada:       models.Inventario{Nombre: ""},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "inventario valido -> sin error y se persiste",
			entrada:       models.Inventario{Nombre: "Laptop Dell", Categoria: "Tecnología", Cantidad: 1},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(inventarioRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearInventario", c.entrada).Return(guardado)
			}
			svc := pi.NewInventarioService(repo)

			creado, err := svc.CrearInventario(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearInventario")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearInventario", c.entrada)
			}
		})
	}
}

func TestInventarioService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(inventarioRepoMock)
	repo.On("BuscarInventarioPorID", 999).Return(models.Inventario{}, false)
	svc := pi.NewInventarioService(repo)

	_, err := svc.BuscarInventario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestInventarioService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(inventarioRepoMock)
	repo.On("BorrarInventario", 999).Return(false)
	svc := pi.NewInventarioService(repo)

	err := svc.BorrarInventario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
