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

// --- Mock Compartido ---
// Doble de prueba del repositorio de acuerdos
type acuerdoRepoMock struct {
	mock.Mock
}

func (m *acuerdoRepoMock) ListarAcuerdos() []models.Acuerdo {
	args := m.Called()
	return args.Get(0).([]models.Acuerdo)
}

func (m *acuerdoRepoMock) BuscarAcuerdoPorID(id int) (models.Acuerdo, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Acuerdo), args.Bool(1)
}

func (m *acuerdoRepoMock) CrearAcuerdo(a models.Acuerdo) models.Acuerdo {
	args := m.Called(a)
	return args.Get(0).(models.Acuerdo)
}

func (m *acuerdoRepoMock) ActualizarAcuerdo(id int, datos models.Acuerdo) (models.Acuerdo, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Acuerdo), args.Bool(1)
}

func (m *acuerdoRepoMock) BorrarAcuerdo(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// Red de seguridad para asegurar que implementa la interfaz del storage de acuerdos
var _ storage.AcuerdoRepository = (*acuerdoRepoMock)(nil)

// --- Tests de Acuerdo ---

// Test con estructura Table Driven para la creación de un acuerdo
func TestAcuerdoService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Acuerdo
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre: "PublicacionID inválido (0) -> ErrVacio o ErrInvalido",
			entrada: models.Acuerdo{
				PublicacionID: 0,
				IDOfertante:   10,
				IDPublicador:  20,
				Tipo:          "Intercambio",
				Estado:        "Pendiente",
			},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre: "Tipo de acuerdo vacío -> ErrVacio",
			entrada: models.Acuerdo{
				PublicacionID: 5,
				IDOfertante:   10,
				IDPublicador:  20,
				Tipo:          "",
				Estado:        "Pendiente",
			},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre: "Acuerdo válido -> sin error y se persiste",
			entrada: models.Acuerdo{
				PublicacionID:            5,
				IDOfertante:              10,
				IDPublicador:             20,
				Tipo:                     "Permuta",
				Estado:                   "Pendiente",
				Mensaje_Inicial:          "Hola, me interesa tu publicación.",
				Confirmacion_Solicitante: false,
			},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// Arrange
			repo := new(acuerdoRepoMock)

			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1
				guardado.CreatedAt = "2026-06-28"
				guardado.UpdatedAt = "2026-06-28"
				// Se usa mock.Anything para evitar fallos por mutaciones menores en el service
				repo.On("CrearAcuerdo", mock.Anything).Return(guardado)
			}

			svc := aiu.NewAcuerdoService(repo)

			// Act
			creado, err := svc.CrearAcuerdo(c.entrada)

			// Assert
			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearAcuerdo", mock.Anything)
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)
				repo.AssertCalled(t, "CrearAcuerdo", mock.Anything)
			}
		})
	}
}

// Test para buscar un acuerdo que no existe
func TestAcuerdoService_Obtener_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(acuerdoRepoMock)
	repo.On("BuscarAcuerdoPorID", 999).Return(models.Acuerdo{}, false)

	svc := aiu.NewAcuerdoService(repo)

	// Act
	_, err := svc.BuscarAcuerdo(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}

// Test para borrar un acuerdo que no existe
func TestAcuerdoService_Borrar_NoEncontrado(t *testing.T) {
	// Arrange
	repo := new(acuerdoRepoMock)
	repo.On("BorrarAcuerdo", 999).Return(false)

	svc := aiu.NewAcuerdoService(repo)

	// Act
	err := svc.BorrarAcuerdo(999)

	// Assert
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
