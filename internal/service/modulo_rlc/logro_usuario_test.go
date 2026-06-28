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
			entrada:       models.Logro_Usuario{FechaDesbl: time.Now()},
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
			svc := rlc.NewLogro_UsuarioService(repo)

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
	svc := rlc.NewLogro_UsuarioService(repo)

	_, err := svc.BuscarLogro_Usuario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

func TestLogro_UsuarioService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(logro_usuarioRepoMock)
	repo.On("BorrarLogro_Usuario", 999).Return(false)
	svc := rlc.NewLogro_UsuarioService(repo)

	err := svc.BorrarLogro_Usuario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
