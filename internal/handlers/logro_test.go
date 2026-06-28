package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/handlers"
	"proyecto-semestral/internal/models"
	rlc "proyecto-semestral/internal/service/modulo_rlc"
	"proyecto-semestral/internal/storage"
)

type logroFake struct {
	porID  map[int]models.Logro
	nextID int
}

func nuevoLogroFake() *logroFake {
	return &logroFake{
		porID:  map[int]models.Logro{},
		nextID: 1,
	}
}

func (f *logroFake) ListarLogro() []models.Logro {
	lista := make([]models.Logro, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *logroFake) BuscarLogroPorID(id int) (models.Logro, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *logroFake) CrearLogro(logro models.Logro) models.Logro {
	logro.ID = f.nextID
	f.nextID++
	f.porID[logro.ID] = logro
	return logro
}

func (f *logroFake) ActualizarLogro(id int, datos models.Logro) (models.Logro, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Logro{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *logroFake) BorrarLogro(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.LogroRepository = (*logroFake)(nil)

// construirEntornoLogro arma un router mínimo solo con la entidad Logro.
func construirEntornoLogro(t *testing.T) http.Handler {
	t.Helper()

	logFake := nuevoLogroFake()
	logService := rlc.NewLogroService(logFake)

	// nil para todos los servicios excepto Logro (5ta posición)
	srv := handlers.NewServer(nil, nil, nil, nil, logService, nil, nil, nil, nil, nil)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/logros", srv.ListarLogro)
		r.Post("/logros", srv.CrearLogro)
	})

	return r
}

func TestCrearLogro_Exitoso(t *testing.T) {
	h := construirEntornoLogro(t)
	body := `{"nombre":"Primer Intercambio","descripcion":"Completaste tu primer trueque","puntos_requeridos":50}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logros", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCrearLogro_NombreVacio(t *testing.T) {
	h := construirEntornoLogro(t)
	body := `{"descripcion":"sin nombre","puntos_requeridos":10}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logros", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListarLogro_Vacio(t *testing.T) {
	h := construirEntornoLogro(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logros", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
