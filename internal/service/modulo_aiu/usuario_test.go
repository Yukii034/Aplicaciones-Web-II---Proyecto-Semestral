package modulo_aiu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	aiu "proyecto-semestral/internal/service/modulo_aiu"
	"proyecto-semestral/internal/storage"
)

// --- Mock ---
// Doble de prueba del repositorio de usuarios
type userRepoMock struct {
	mock.Mock
}

func (m *userRepoMock) ListarUsuarios() []models.Usuario {
	args := m.Called()
	return args.Get(0).([]models.Usuario)
}

func (m *userRepoMock) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Usuario), args.Bool(1)
}

func (m *userRepoMock) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	args := m.Called(email)
	return args.Get(0).(models.Usuario), args.Bool(1)
}

func (m *userRepoMock) CrearUsuario(u models.Usuario) models.Usuario {
	args := m.Called(u)
	return args.Get(0).(models.Usuario)
}

func (m *userRepoMock) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Usuario), args.Bool(1)
}

func (m *userRepoMock) BorrarUsuario(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// Red de seguridad para asegurar que implementa la interfaz del storage
var _ storage.UserRepository = (*userRepoMock)(nil)

// --- Tests ---

// Test con estructura Table Driven para la creación de un usuario
func TestUsuarioService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Usuario
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "nombre vacío -> ErrVacio",
			entrada:       models.Usuario{Nombre: "", Email: "juan@mail.com", PasswordHash: "hash123", Tipo: "cliente"},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "email vacío -> ErrVacio",
			entrada:       models.Usuario{Nombre: "Juan", Email: "", PasswordHash: "hash123", Tipo: "cliente"},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "usuario válido -> sin error y se persiste",
			entrada:       models.Usuario{Nombre: "Juan Pérez", Email: "juan@mail.com", PasswordHash: "hash123", Tipo: "Administrador", Ciudad: "Guayaquil", Telefono: "0999999999", Reputacion: "Buena"},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			repo := new(userRepoMock)

			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1 // Simula el ID autoincremental que genera GORM
				repo.On("CrearUsuario", c.entrada).Return(guardado)
			}

			svc := aiu.NewUsuarioService(repo)

			// Act
			creado, err := svc.CrearUsuario(c.entrada)

			// Assert
			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearUsuario")
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearUsuario", c.entrada)
			}
		})
	}
}

// Test para buscar un usuario que no existe
func TestUsuarioService_Obtener_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(userRepoMock)
	repo.On("BuscarUsuarioPorID", 999).Return(models.Usuario{}, false)

	svc := aiu.NewUsuarioService(repo)

	// Act
	_, err := svc.BuscarUsuario(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

// Test para borrar un usuario que no existe
func TestUsuarioService_Borrar_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(userRepoMock)
	repo.On("BorrarUsuario", 999).Return(false)

	svc := aiu.NewUsuarioService(repo)

	// Act
	err := svc.BorrarUsuario(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
