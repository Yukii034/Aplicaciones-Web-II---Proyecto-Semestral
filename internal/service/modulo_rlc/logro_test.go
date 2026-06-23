package modulo_rlc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	rlc "proyecto-semestral/internal/service/modulo_rlc"
	"proyecto-semestral/internal/storage"
)

// --- Mock ---
type logroRepoMock struct {
	mock.Mock
}

func (m *logroRepoMock) ListarLogro() []models.Logro {
	args := m.Called()
	return args.Get(0).([]models.Logro)
}

func (m *logroRepoMock) BuscarLogroPorID(id int) (models.Logro, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Logro), args.Bool(1)
}

func (m *logroRepoMock) CrearLogro(lu models.Logro) models.Logro {
	args := m.Called(lu)
	return args.Get(0).(models.Logro)
}

func (m *logroRepoMock) ActualizarLogro(id int, datos models.Logro) (models.Logro, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Logro), args.Bool(1)
}

func (m *logroRepoMock) BorrarLogro(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// red de seguridad
var _ storage.LogroRepository = (*logroRepoMock)(nil)

// --- Tests ---
func TestLogroService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Logro
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "nombre vacio -> ErrVacio",
			entrada:       models.Logro{},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "logro valido -> sin error y se persiste",
			entrada:       models.Logro{Nombre: "Primer intercambio", Descripcion: "Lograste tu primer intercambio", PuntosRequeridos: 20},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(logroRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearLogro", c.entrada).Return(guardado)
			}
			svc := rlc.NewLogroService(repo)

			creado, err := svc.CrearLogro(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearLogro")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearLogro", c.entrada)
			}
		})
	}
}

func TestLogro_Obtener_NoEncontrado(t *testing.T) {
	repo := new(logroRepoMock)
	repo.On("BuscarLogroPorID", 999).Return(models.Logro{}, false)
	svc := rlc.NewLogroService(repo)

	_, err := svc.BuscarLogro(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestLogroService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(logroRepoMock)
	repo.On("BorrarLogro", 999).Return(false)
	svc := rlc.NewLogroService(repo)

	err := svc.BorrarLogro(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
