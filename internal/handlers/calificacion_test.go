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

type CalificacionFake struct {
	porID  map[int]models.Calificacion
	nextID int
}

func nuevoCalificacionFake() *CalificacionFake {
	return &CalificacionFake{
		porID:  map[int]models.Calificacion{},
		nextID: 1,
	}
}

func (f *CalificacionFake) ListarCalificacion() []models.Calificacion {
	lista := make([]models.Calificacion, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *CalificacionFake) BuscarCalificacionPorID(id int) (models.Calificacion, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *CalificacionFake) CrearCalificacion(ca models.Calificacion) models.Calificacion {
	ca.ID = f.nextID
	f.nextID++
	f.porID[ca.ID] = ca
	return ca
}

func (f *CalificacionFake) ActualizarCalificacion(id int, datos models.Calificacion) (models.Calificacion, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Calificacion{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *CalificacionFake) BorrarCalificacion(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.CalificacionRepository = (*CalificacionFake)(nil)

func construirEntornoCalificacion(t *testing.T) http.Handler {
	t.Helper()

	caFake := nuevoCalificacionFake()
	caService := rlc.NewCalificacionService(caFake)

	srv := handlers.NewServer(handlers.Deps{Calificacion: caService})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/calificacion", srv.ListarCalificacion)
		r.Post("/calificacion", srv.CrearCalificacion)
		r.Get("/calificacion/{id}", srv.ObtenerCalificacion)
		r.Put("/calificacion/{id}", srv.ActualizarCalificacion)
		r.Delete("/calificacion/{id}", srv.BorrarCalificacion)
	})

	return r
}

func TestCrearCalificacion_Exitoso(t *testing.T) {
	h := construirEntornoCalificacion(t)
	body := `{"comentarios":"Excelente intercambio","usuario_id":1,"acuerdo_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calificacion", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCrearCalificacion_ComentarioVacio(t *testing.T) {
	h := construirEntornoCalificacion(t)
	body := `{"comentarios":"","usuario_id":1,"acuerdo_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calificacion", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListarCalificacion_Vacio(t *testing.T) {
	h := construirEntornoCalificacion(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calificacion", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestObtenerCalificacion_NoExiste(t *testing.T) {
	h := construirEntornoCalificacion(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calificacion/999", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCrearCalificacion_SinUsuarioID(t *testing.T) {
	h := construirEntornoCalificacion(t)
	body := `{"comentarios":"Buen intercambio","acuerdo_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calificacion", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
