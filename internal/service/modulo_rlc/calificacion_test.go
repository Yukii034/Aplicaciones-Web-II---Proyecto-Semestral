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
type calificacionRepoMock struct {
	mock.Mock
}

func (m *calificacionRepoMock) ListarCalificacion() []models.Calificacion {
	args := m.Called()
	return args.Get(0).([]models.Calificacion)
}

func (m *calificacionRepoMock) BuscarCalificacionPorID(id int) (models.Calificacion, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Calificacion), args.Bool(1)
}

func (m *calificacionRepoMock) CrearCalificacion(lu models.Calificacion) models.Calificacion {
	args := m.Called(lu)
	return args.Get(0).(models.Calificacion)
}

func (m *calificacionRepoMock) ActualizarCalificacion(id int, datos models.Calificacion) (models.Calificacion, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Calificacion), args.Bool(1)
}

func (m *calificacionRepoMock) BorrarCalificacion(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// red de seguridad
var _ storage.CalificacionRepository = (*calificacionRepoMock)(nil)

// --- Tests ---
func TestCalificacionService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Calificacion
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "comentario vacio -> ErrVacio",
			entrada:       models.Calificacion{},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "comentario valido -> sin error y se persiste",
			entrada:       models.Calificacion{Comentarios: "Excelete intercambio"},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(calificacionRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearCalificacion", c.entrada).Return(guardado)
			}
			svc := rlc.NewCalificacionService(repo)

			creado, err := svc.CrearCalificacion(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearCalificacion")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearCalificacion", c.entrada)
			}
		})
	}
}

func TestCalificacion_Obtener_NoEncontrado(t *testing.T) {
	repo := new(calificacionRepoMock)
	repo.On("BuscarCalificacionPorID", 999).Return(models.Calificacion{}, false)
	svc := rlc.NewCalificacionService(repo)

	_, err := svc.BuscarCalificacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestCalificacionService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(calificacionRepoMock)
	repo.On("BorrarCalificacion", 999).Return(false)
	svc := rlc.NewCalificacionService(repo)

	err := svc.BorrarCalificacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
