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

// usuario
type usuarioRepoMockCal struct {
	mock.Mock
}

func (m *usuarioRepoMockCal) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Usuario), args.Bool(1)
}

func (m *usuarioRepoMockCal) ListarUsuarios() []models.Usuario { return nil }
func (m *usuarioRepoMockCal) BuscarUsuarioPorEmail(e string) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *usuarioRepoMockCal) CrearUsuario(u models.Usuario) models.Usuario { return models.Usuario{} }
func (m *usuarioRepoMockCal) ActualizarUsuario(id int, d models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *usuarioRepoMockCal) BorrarUsuario(id int) bool { return false }

var _ storage.UserRepository = (*usuarioRepoMockCal)(nil)

// Acuerdo
type acuerdoRepoMockCal struct {
	mock.Mock
}

func (m *acuerdoRepoMockCal) BuscarAcuerdoPorID(id int) (models.Acuerdo, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Acuerdo), args.Bool(1)
}

func (m *acuerdoRepoMockCal) ListarAcuerdos() []models.Acuerdo             { return nil }
func (m *acuerdoRepoMockCal) CrearAcuerdo(a models.Acuerdo) models.Acuerdo { return models.Acuerdo{} }
func (m *acuerdoRepoMockCal) ActualizarAcuerdo(id int, a models.Acuerdo) (models.Acuerdo, bool) {
	return models.Acuerdo{}, false
}
func (m *acuerdoRepoMockCal) BorrarAcuerdo(id int) bool { return false }

var _ storage.AcuerdoRepository = (*acuerdoRepoMockCal)(nil)

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
			entrada:       models.Calificacion{Comentarios: "Excelente intercambio", UsuarioID: 1, AcuerdoID: 1},
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
			usuRepo := new(usuarioRepoMockCal)
			acuRepo := new(acuerdoRepoMockCal)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearCalificacion", c.entrada).Return(guardado)
				usuRepo.On("BuscarUsuarioPorID", 1).Return(models.Usuario{ID: 1}, true)
				acuRepo.On("BuscarAcuerdoPorID", 1).Return(models.Acuerdo{ID: 1}, true)
			}
			svc := rlc.NewCalificacionService(repo, usuRepo, acuRepo)

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

func TestCalificacionService_Obtener_NoEncontrado(t *testing.T) {
	repo := new(calificacionRepoMock)
	repo.On("BuscarCalificacionPorID", 999).Return(models.Calificacion{}, false)
	usrRepo := new(usuarioRepoMockCal)
	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	_, err := svc.BuscarCalificacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestCalificacionService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(calificacionRepoMock)
	repo.On("BorrarCalificacion", 999).Return(false)
	usrRepo := new(usuarioRepoMockCal)
	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	err := svc.BorrarCalificacion(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestCalificacionService_Actualizar_NoEncontrado(t *testing.T) {
	repo := new(calificacionRepoMock)
	usrRepo := new(usuarioRepoMockCal)

	usrRepo.On("BuscarUsuarioPorID", 0).Return(models.Usuario{ID: 0}, true)

	datos := models.Calificacion{Comentarios: "Buen intercambio", UsuarioID: 1}
	repo.On("ActualizarCalificacion", 999, datos).Return(models.Calificacion{}, false)

	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	_, err := svc.ActualizarCalificacion(999, datos)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestCalificacionService_Actualizar_ComentarioVacio(t *testing.T) {
	repo := new(calificacionRepoMock)
	usrRepo := new(usuarioRepoMockCal)
	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	_, err := svc.ActualizarCalificacion(1, models.Calificacion{Comentarios: ""})

	require.ErrorIs(t, err, service.ErrVacio)
	repo.AssertNotCalled(t, "ActualizarCalificacion")
}

func TestCalificacionService_Listar(t *testing.T) {
	repo := new(calificacionRepoMock)
	esperado := []models.Calificacion{
		{ID: 1, Comentarios: "muy buen intercambio"},
		{ID: 2, Comentarios: "muy buena donacion"},
	}
	repo.On("ListarCalificacion").Return(esperado)
	usrRepo := new(usuarioRepoMockCal)
	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	lista := svc.ListarCalificacion()

	assert.Len(t, lista, 2)
	assert.Equal(t, "muy buen intercambio", lista[0].Comentarios)
	repo.AssertExpectations(t)
}

func TestCalificacionService_Obtener_Exitoso(t *testing.T) {
	repo := new(calificacionRepoMock)
	esperado := models.Calificacion{ID: 1, Comentarios: "muy amable"}
	repo.On("BuscarCalificacionPorID", 1).Return(esperado, true)
	usrRepo := new(usuarioRepoMockCal)
	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	encontrada, err := svc.BuscarCalificacion(1)

	require.NoError(t, err)
	assert.Equal(t, "muy amable", encontrada.Comentarios)
	repo.AssertExpectations(t)
}

func TestCalificacionService_Borrar_Exitoso(t *testing.T) {
	repo := new(calificacionRepoMock)
	repo.On("BorrarCalificacion", 1).Return(true)
	usrRepo := new(usuarioRepoMockCal)
	svc := rlc.NewCalificacionService(repo, usrRepo, new(acuerdoRepoMockCal))

	err := svc.BorrarCalificacion(1)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
