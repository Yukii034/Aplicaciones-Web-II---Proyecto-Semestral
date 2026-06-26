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
// doble de prueba del repositorio de inventario
type inventarioRepoMock struct {
	mock.Mock
}

// metodos de la interfaz unica de inventario, muestran qué necesitan que se mande a la función y qué es lo que devuelven

func (m *inventarioRepoMock) ListarInventario() []models.Inventario {
	args := m.Called()
	return args.Get(0).([]models.Inventario)
}

func (m *inventarioRepoMock) BuscarInventarioPorID(id int) (models.Inventario, bool) {
	args := m.Called(id)
	return args.Get(0).(models.Inventario), args.Bool(1)
}

func (m *inventarioRepoMock) CrearInventario(i models.Inventario) models.Inventario {
	args := m.Called(i)
	return args.Get(0).(models.Inventario)
}

func (m *inventarioRepoMock) ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool) {
	args := m.Called(id, datos)
	return args.Get(0).(models.Inventario), args.Bool(1)
}

func (m *inventarioRepoMock) BorrarInventario(id int) bool {
	args := m.Called(id)
	return args.Bool(0)
}

// red de seguridad, para que cumpla el contrato y pueda compilar
var _ storage.InventarioRepository = (*inventarioRepoMock)(nil)

// --- Tests ---
// estructura Table Driven, tiene algunos casos para recorrer
func TestInventarioService_Crear(t *testing.T) {
	casos := []struct {
		nombre        string
		entrada       models.Inventario
		errEsperado   error
		debePersistir bool
	}{
		{
			nombre:        "nombre vacio -> ErrVacio", // si no hay nombre debe de salir el error de que está vacío
			entrada:       models.Inventario{Nombre: ""},
			errEsperado:   service.ErrVacio,
			debePersistir: false,
		},
		{
			nombre:        "inventario valido -> sin error y se persiste", // datos validos ingresados, se espera que salga sin errores
			entrada:       models.Inventario{Nombre: "Laptop Dell", Categoria: "Tecnología", Cantidad: 1},
			errEsperado:   nil,
			debePersistir: true,
		},
	}

	// recorre los casos
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) { // t.run crea un subtest por cada caso
			// Arrange - prepara el mock nuevo
			repo := new(inventarioRepoMock)
			// se guarda el inventario con un id asignado
			if c.debePersistir {
				guardado := c.entrada
				guardado.ID = 1                                        // le asigna un id como gorm, ya que esta es una copia sin este
				repo.On("CrearInventario", c.entrada).Return(guardado) // solo se prepara el mock si el repo es llamado
			}
			svc := pi.NewInventarioService(repo)

			// Act - ejecuta, creando el nuevo inventario
			creado, err := svc.CrearInventario(c.entrada)

			// Assert - verifica que se haya creado el nuevo inventario
			if c.errEsperado != nil {
				require.ErrorIs(t, err, c.errEsperado)
				repo.AssertNotCalled(t, "CrearInventario") // si hay errores no llama nunca toca el repo, caso con error
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, creado.ID)                      // se verifica si es el dato guardado con id 1
				repo.AssertCalled(t, "CrearInventario", c.entrada) // saber si llamó al metodo, caso exitoso
			}
		})
	}
}

func TestInventarioService_Obtener_NoEncontrado(t *testing.T) {
	// Arrange - prepara el mock nuevo - un struct vacio
	repo := new(inventarioRepoMock)

	// Act - Cuando alguien llame BuscarInventarioPorID con el argumento 999, devuelve un Inventario vacío y false
	// On - registra qué método esperas que se llame y con qué argumento
	// Return - define qué devuelve cuando ese método se llama
	repo.On("BuscarInventarioPorID", 999).Return(models.Inventario{}, false)

	// Crea el service de inventario pero le inyecta el mock en vez del real
	svc := pi.NewInventarioService(repo)

	// Llama al metodo real del service con el id inexistente
	// _ descarta el primer valor porque solo se necesita saber qué error devuelve
	_, err := svc.BuscarInventario(999)

	// Assert - verifica que el error sea "ErrNoEncontrado" del service
	// con require, si falla para el test, con assert seguiría corriendo el test en busca de otros errores
	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t) // para saber si llamó al método correcto
}

func TestInventarioService_Borrar_NoEncontrado(t *testing.T) {
	repo := new(inventarioRepoMock)
	repo.On("BorrarInventario", 999).Return(false)
	svc := pi.NewInventarioService(repo)

	err := svc.BorrarInventario(999)

	require.ErrorIs(t, err, service.ErrNoEncontrado)
	repo.AssertExpectations(t)
}
