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
type reputacionRepoMock struct {
	mock.Mock
}

func (m *reputacionRepoMock) ListarReputacion() []models.Reputacion {
	args := m.Called()
	return args.Get(0).([]models.Reputacion)
}

func (m *reputacionRepoMock) BuscarReputacionPorID(id int) (models.Reputacion, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Reputacion), args.Bool(1)
}

func (m *reputacionRepoMock) CrearReputacion(r models.Reputacion) models.Reputacion {
	args := m.Called(r)
	return args.Get(0).(models.Reputacion)
}

func (m *reputacionRepoMock) ActualizarReputacion(id int, datos models.Reputacion) (models.Reputacion, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Reputacion), args.Bool(1)
}

func (m *reputacionRepoMock) BorrarReputacion(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// red de seguridad
var _ storage.ReputacionRepository = (*reputacionRepoMock)(nil)

// --- Tests ---
func TestReputacionService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Reputacion
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "puntos totales vacio -> ErrVacio",
			entrada:       models.Reputacion{PuntosTotales: 0},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "reputacion valido -> sin error y se persiste",
			entrada:       models.Reputacion{PuntosTotales: 100, Nivel: 20, AcuerdosCompl: 20, CalificacionPromedio: 4.5},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(reputacionRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearReputacion", c.entrada).Return(guardado)
			}
			svc := rlc.NewReputacionService(repo)

			creado, err := svc.CrearReputacion(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearReputacion")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearReputacion", c.entrada)
			}
		})
	}
}

func TestReputacionService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(reputacionRepoMock)
	repo.On("BuscarReputacionPorID", 999).Return(models.Reputacion{}, false)
	svc := rlc.NewReputacionService(repo)

	_, err := svc.BuscarReputacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestReputacionService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(reputacionRepoMock)
	repo.On("BorrarReputacion", 999).Return(false)
	svc := rlc.NewReputacionService(repo)

	err := svc.BorrarReputacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
