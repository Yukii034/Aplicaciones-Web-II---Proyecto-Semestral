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
			entrada:       models.Reputacion{PuntosTotales: 100, Nivel: 2, AcuerdosCompl: 3, CalificacionPromedio: 4.5, UsuarioID: 1},
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

func TestReputacionService_Actualizar_NoEncontrado(t *testing.T) {
	repo := new(reputacionRepoMock)
	repo.On("ActualizarReputacion", 999, models.Reputacion{PuntosTotales: 150, Nivel: 2, AcuerdosCompl: 4, CalificacionPromedio: 4.5, UsuarioID: 1}).Return(models.Reputacion{}, false)
	svc := rlc.NewReputacionService(repo)

	_, err := svc.ActualizarReputacion(999, models.Reputacion{PuntosTotales: 150, Nivel: 2, AcuerdosCompl: 4, CalificacionPromedio: 4.5, UsuarioID: 1})

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestReputacionService_Actualizar_Puntos_TotalesVacio(t *testing.T) {
	repo := new(reputacionRepoMock)
	svc := rlc.NewReputacionService(repo)

	_, err := svc.ActualizarReputacion(1, models.Reputacion{PuntosTotales: 0})

	require.ErrorIs(t, err, service.ErrVacio)
	repo.AssertNotCalled(t, "ActualizarReputacion")
}

func TestReputacionService_Listar(t *testing.T) {
	repo := new(reputacionRepoMock)
	esperado := []models.Reputacion{
		{ID: 1, PuntosTotales: 150, Nivel: 2, AcuerdosCompl: 4, CalificacionPromedio: 4.5},
		{ID: 2, PuntosTotales: 200, Nivel: 3, AcuerdosCompl: 7, CalificacionPromedio: 4.7},
	}
	repo.On("ListarReputacion").Return(esperado)
	svc := rlc.NewReputacionService(repo)

	lista := svc.ListarReputacion()

	assert.Len(t, lista, 2)
	assert.Equal(t, 150, lista[0].PuntosTotales)
	repo.AssertExpectations(t)
}

func TestReputacionService_Obtener_Exitoso(t *testing.T) {
	repo := new(reputacionRepoMock)
	esperado := models.Reputacion{ID: 1, PuntosTotales: 150, Nivel: 2, AcuerdosCompl: 4, CalificacionPromedio: 4.5}
	repo.On("BuscarReputacionPorID", 1).Return(esperado, true)
	svc := rlc.NewReputacionService(repo)

	encontrado, err := svc.BuscarReputacion(1)

	require.NoError(t, err)
	assert.Equal(t, 150, encontrado.PuntosTotales)
	repo.AssertExpectations(t)
}

func TestReputacionService_Borrar_Exitoso(t *testing.T) {
	repo := new(reputacionRepoMock)
	repo.On("BorrarReputacion", 1).Return(true)
	svc := rlc.NewReputacionService(repo)

	err := svc.BorrarReputacion(1)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
