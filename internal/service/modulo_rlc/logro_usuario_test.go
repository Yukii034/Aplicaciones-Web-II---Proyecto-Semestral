package modulo_rlc_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	rlc "proyecto-semestral/internal/service/modulo_rlc"
	"proyecto-semestral/internal/storage"
)

// --- Mock ---
type logro_usuarioRepoMock struct {
	mock.Mock
}

func (m *logro_usuarioRepoMock) ListarLogro_Usuario() []models.Logro_Usuario {
	args := m.Called()
	return args.Get(0).([]models.Logro_Usuario)
}

func (m *logro_usuarioRepoMock) BuscarLogro_UsuarioPorID(id int) (models.Logro_Usuario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Logro_Usuario), args.Bool(1)
}

func (m *logro_usuarioRepoMock) CrearLogro_Usuario(lu models.Logro_Usuario) models.Logro_Usuario {
	args := m.Called(lu)
	return args.Get(0).(models.Logro_Usuario)
}

func (m *logro_usuarioRepoMock) ActualizarLogro_Usuario(id int, datos models.Logro_Usuario) (models.Logro_Usuario, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Logro_Usuario), args.Bool(1)
}

func (m *logro_usuarioRepoMock) BorrarLogro_Usuario(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// red de seguridad
var _ storage.Logro_UsuarioRepository = (*logro_usuarioRepoMock)(nil)

// usuario
type usuarioRepoMockCalll struct {
	mock.Mock
}

func (m *usuarioRepoMockCalll) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Usuario), args.Bool(1)
}

func (m *usuarioRepoMockCalll) ListarUsuarios() []models.Usuario { return nil }
func (m *usuarioRepoMockCalll) BuscarUsuarioPorEmail(e string) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *usuarioRepoMockCalll) CrearUsuario(u models.Usuario) models.Usuario { return models.Usuario{} }
func (m *usuarioRepoMockCalll) ActualizarUsuario(id int, d models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}
func (m *usuarioRepoMockCalll) BorrarUsuario(id int) bool { return false }

var _ storage.UserRepository = (*usuarioRepoMockCalll)(nil)

// logro
type LogroRepoMockCal struct {
	mock.Mock
}

func (m *LogroRepoMockCal) BuscarLogroPorID(id int) (models.Logro, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Logro), args.Bool(1)
}

func (m *LogroRepoMockCal) ListarLogro() []models.Logro { return nil }
func (m *LogroRepoMockCal) BuscarLogroPorEmail(e string) (models.Logro, bool) {
	return models.Logro{}, false
}
func (m *LogroRepoMockCal) CrearLogro(u models.Logro) models.Logro { return models.Logro{} }
func (m *LogroRepoMockCal) ActualizarLogro(id int, d models.Logro) (models.Logro, bool) {
	return models.Logro{}, false
}
func (m *LogroRepoMockCal) BorrarLogro(id int) bool { return false }

var _ storage.LogroRepository = (*LogroRepoMockCal)(nil)

// --- Tests ---
func TestLogro_UsuarioService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Logro_Usuario
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "fecha desbloqueada vacio -> ErrVacio",
			entrada:       models.Logro_Usuario{},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "logro usuario valido -> sin error y se persiste",
			entrada:       models.Logro_Usuario{FechaDesbl: time.Now(), UsuarioID: 1, LogroID: 1},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := new(logro_usuarioRepoMock)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearLogro_Usuario", c.entrada).Return(guardado)
			}
			usrRepo := new(usuarioRepoMockCalll)
			logRepo := new(LogroRepoMockCal)
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				repo.On("CrearLogro_Usuario", c.entrada).Return(guardado)
				usrRepo.On("BuscarUsuarioPorID", 1).Return(models.Usuario{ID: 1}, true)
				logRepo.On("BuscarLogroPorID", 1).Return(models.Logro{ID: 1}, true)
			}
			svc := rlc.NewLogro_UsuarioService(repo, logRepo, usrRepo)

			creado, err := svc.CrearLogro_Usuario(c.entrada)

			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearLogro_usuario")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearLogro_Usuario", c.entrada)
			}
		})
	}
}

func TestLogro_Usuario_Obtener_NoEncontrado(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	repo.On("BuscarLogro_UsuarioPorID", 999).Return(models.Logro_Usuario{}, false)
	logRepo := new(LogroRepoMockCal)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	_, err := svc.BuscarLogro_Usuario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestLogro_UsuarioService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	repo.On("BorrarLogro_Usuario", 999).Return(false)
	logRepo := new(LogroRepoMockCal)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	err := svc.BorrarLogro_Usuario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestLogro_UsuarioService_Actualizar_NoEncontrado(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	logRepo := new(LogroRepoMockCal)
	fecha := time.Now()
	logRepo.On("BuscarLogroPorID")
	datos := models.Logro_Usuario{FechaDesbl: fecha, UsuarioID: 1, LogroID: 1}
	repo.On("ActualizarLogro_Usuario", 999, datos).Return(models.Logro_Usuario{}, false)
	logRepo.On("BuscarLogroPorID", 0).Return(models.Logro{ID: 0}, true)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	_, err := svc.ActualizarLogro_Usuario(999, datos)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestLogro_UsuarioService_Actualizar_FechaVacia(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	logRepo := new(LogroRepoMockCal)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	_, err := svc.ActualizarLogro_Usuario(1, models.Logro_Usuario{})

	require.ErrorIs(t, err, service.ErrVacio)
	repo.AssertNotCalled(t, "ActualizarLogro_Usuario")
}

func TestLogro_UsuarioService_Listar(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	fecha := time.Now()
	esperado := []models.Logro_Usuario{
		{ID: 1, FechaDesbl: fecha, UsuarioID: 1, LogroID: 1},
		{ID: 2, FechaDesbl: fecha, UsuarioID: 2, LogroID: 2},
	}
	repo.On("ListarLogro_Usuario").Return(esperado)
	logRepo := new(LogroRepoMockCal)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	lista := svc.ListarLogro_Usuario()

	assert.Len(t, lista, 2)
	assert.Equal(t, 1, lista[0].UsuarioID)
	repo.AssertExpectations(t)
}

func TestLogro_UsuarioService_Obtener_Exitoso(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	fecha := time.Now()
	esperado := models.Logro_Usuario{ID: 1, FechaDesbl: fecha, UsuarioID: 1, LogroID: 1}
	repo.On("BuscarLogro_UsuarioPorID", 1).Return(esperado, true)
	logRepo := new(LogroRepoMockCal)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	encontrado, err := svc.BuscarLogro_Usuario(1)

	require.NoError(t, err)
	assert.Equal(t, 1, encontrado.UsuarioID)
	repo.AssertExpectations(t)
}

func TestLogro_UsuarioService_Borrar_Exitoso(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	repo.On("BorrarLogro_Usuario", 1).Return(true)
	logRepo := new(LogroRepoMockCal)
	svc := rlc.NewLogro_UsuarioService(repo, logRepo, new(usuarioRepoMockCalll))

	err := svc.BorrarLogro_Usuario(1)

	require.NoError(t, err)
	repo.AssertExpectations(t)
}
