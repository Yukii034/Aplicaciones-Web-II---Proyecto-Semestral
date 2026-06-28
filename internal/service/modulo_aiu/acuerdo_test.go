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

// CORRECCIÓN AQUÍ: Flexibilidad en el Mock para aceptar tanto estructuras como punteros si fuera necesario
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

var _ storage.AcuerdoRepository = (*acuerdoRepoMock)(nil)

// --- Tests de Acuerdo ---

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

				// CORRECCIÓN: Usamos MatchedBy para que acepte cualquier estructura models.Acuerdo
				// sin importar si el service alteró algún campo interno (como fechas o estados por defecto).
				repo.On("CrearAcuerdo", mock.MatchedBy(func(_ models.Acuerdo) bool { return true })).Return(guardado)
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
				repo.AssertExpectations(t) // Asegura que se llamó al mock esperado correctamente
			}
		})
	}
}
