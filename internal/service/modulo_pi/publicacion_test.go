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
type publicacionRepoMock struct {
	mock.Mock
}

func (m *publicacionRepoMock) ListarPublicacion() []models.Publicacion {
	args := m.Called()
	return args.Get(0).([]models.Publicacion)
}

func (m *publicacionRepoMock) BuscarPublicacionPorID(id int) (models.Publicacion, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Publicacion), args.Bool(1)
}

func (m *publicacionRepoMock) CrearPublicacion(p models.Publicacion) models.Publicacion {
	args := m.Called(p)
	return args.Get(0).(models.Publicacion)
}

func (m *publicacionRepoMock) ActualizarPublicacion(id int, datos models.Publicacion) (models.Publicacion, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Publicacion), args.Bool(1)
}

func (m *publicacionRepoMock) BorrarPublicacion(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

var _ storage.PublicacionRepository = (*publicacionRepoMock)(nil)

// --- Tests ---
func TestPublicacionService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Publicacion
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "titulo vacio -> ErrVacio",
			entrada:       models.Publicacion{Titulo: ""},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "publicacion valida -> sin error y se persiste",
			entrada:       models.Publicacion{Titulo: "Cambio laptop", TipoOferta: "intercambio"},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(publicacionRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearPublicacion", c.entrada).Return(guardado)
			}
			svc := pi.NewPublicacionService(repo)

			creada, err := svc.CrearPublicacion(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearPublicacion")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creada.ID)
				repo.AssertCalled(t, "CrearPublicacion", c.entrada)
			}
		})
	}
}

func TestPublicacionService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(publicacionRepoMock)
	repo.On("BuscarPublicacionPorID", 999).Return(models.Publicacion{}, false)
	svc := pi.NewPublicacionService(repo)

	_, err := svc.BuscarPublicacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(publicacionRepoMock)
	repo.On("BorrarPublicacion", 999).Return(false)
	svc := pi.NewPublicacionService(repo)

	err := svc.BorrarPublicacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
