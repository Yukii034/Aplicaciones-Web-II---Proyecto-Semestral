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

// Mock publicacion
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

// Mock inventario (para relacion)
type inventarioRepoMockP struct {
	mock.Mock
}

func (m *inventarioRepoMockP) ListarInventario() []models.Inventario {
	args := m.Called()
	return args.Get(0).([]models.Inventario)
}

func (m *inventarioRepoMockP) BuscarInventarioPorID(id int) (models.Inventario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Inventario), args.Bool(1)
}

func (m *inventarioRepoMockP) CrearInventario(i models.Inventario) models.Inventario {
	args := m.Called(i)
	return args.Get(0).(models.Inventario)
}

func (m *inventarioRepoMockP) ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Inventario), args.Bool(1)
}

func (m *inventarioRepoMockP) BorrarInventario(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

var _ storage.InventarioRepository = (*inventarioRepoMockP)(nil)

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
			entrada:       models.Publicacion{Titulo: "Cambio laptop", TipoOferta: "intercambio", InventarioID: 1},
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
			invRepo := new(inventarioRepoMockP)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearPublicacion", c.entrada).Return(guardado)
				// el inventario con ID 1 existe
				invRepo.On("BuscarInventarioPorID", 1).Return(models.Inventario{ID: 1}, true)
			}

			svc := pi.NewPublicacionService(repo, invRepo)

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
	invRepo := new(inventarioRepoMockP)
	svc := pi.NewPublicacionService(repo, invRepo)

	_, err := svc.BuscarPublicacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(publicacionRepoMock)
	repo.On("BorrarPublicacion", 999).Return(false)
	invRepo := new(inventarioRepoMockP)
	svc := pi.NewPublicacionService(repo, invRepo)

	err := svc.BorrarPublicacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Actualizar_NoEncontrado(t *testing.T) {
	repo := new(publicacionRepoMock)
	invRepo := new(inventarioRepoMockP)

	// el inventario sí existe (para que pase esa validación)
	invRepo.On("BuscarInventarioPorID", 0).Return(models.Inventario{ID: 0}, true)

	// pero la publicacion con id 999 no existe
	repo.On("ActualizarPublicacion", 999, models.Publicacion{Titulo: "Cambio laptop", TipoOferta: "intercambio"}).Return(models.Publicacion{}, false)

	svc := pi.NewPublicacionService(repo, invRepo)

	_, err := svc.ActualizarPublicacion(999, models.Publicacion{Titulo: "Cambio laptop", TipoOferta: "intercambio"})

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Actualizar_TituloVacio(t *testing.T) {
	repo := new(publicacionRepoMock)
	invRepo := new(inventarioRepoMockP)
	svc := pi.NewPublicacionService(repo, invRepo)

	_, err := svc.ActualizarPublicacion(1, models.Publicacion{Titulo: ""})

	require.ErrorIs(t, err, service.ErrVacio)
	repo.AssertNotCalled(t, "ActualizarPublicacion")
}

func TestPublicacionService_Listar(t *testing.T) {
	repo := new(publicacionRepoMock)
	esperado := []models.Publicacion{
		{ID: 1, Titulo: "Cambio laptop", TipoOferta: "intercambio"},
		{ID: 2, Titulo: "Dono microondas", TipoOferta: "donacion"},
	}
	repo.On("ListarPublicacion").Return(esperado)
	invRepo := new(inventarioRepoMockP)
	svc := pi.NewPublicacionService(repo, invRepo)

	lista := svc.ListarPublicacion()

	assert.Len(t, lista, 2)
	assert.Equal(t, "Cambio laptop", lista[0].Titulo)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Obtener_Exitoso(t *testing.T) {
	repo := new(publicacionRepoMock)
	esperado := models.Publicacion{ID: 1, Titulo: "Cambio laptop", TipoOferta: "intercambio"}
	repo.On("BuscarPublicacionPorID", 1).Return(esperado, true)
	invRepo := new(inventarioRepoMockP)
	svc := pi.NewPublicacionService(repo, invRepo)

	encontrada, err := svc.BuscarPublicacion(1)

	require.NoError(t, err)
	assert.Equal(t, "Cambio laptop", encontrada.Titulo)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Borrar_Exitoso(t *testing.T) {
	repo := new(publicacionRepoMock)
	repo.On("BorrarPublicacion", 1).Return(true)
	invRepo := new(inventarioRepoMockP)
	svc := pi.NewPublicacionService(repo, invRepo)

	err := svc.BorrarPublicacion(1)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPublicacionService_Crear_InventarioNoExiste(t *testing.T) {
	repo := new(publicacionRepoMock)
	invRepo := new(inventarioRepoMockP)
	invRepo.On("BuscarInventarioPorID", 99).Return(models.Inventario{}, false)
	svc := pi.NewPublicacionService(repo, invRepo)

	_, err := svc.CrearPublicacion(models.Publicacion{
		Titulo:       "Cambio laptop",
		TipoOferta:   "intercambio",
		InventarioID: 99,
	})

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertNotCalled(t, "CrearPublicacion")
}

func TestPublicacionService_Actualizar_InventarioNoExiste(t *testing.T) {
	repo := new(publicacionRepoMock)
	invRepo := new(inventarioRepoMockP)
	invRepo.On("BuscarInventarioPorID", 99).Return(models.Inventario{}, false)
	svc := pi.NewPublicacionService(repo, invRepo)

	_, err := svc.ActualizarPublicacion(1, models.Publicacion{
		Titulo:       "Cambio laptop",
		TipoOferta:   "intercambio",
		InventarioID: 99,
	})

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertNotCalled(t, "ActualizarPublicacion")
}
