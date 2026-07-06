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

type ReputacionFake struct {
	porID  map[int]models.Reputacion
	nextID int
}

func nuevoReputacionFake() *ReputacionFake {
	return &ReputacionFake{
		porID:  map[int]models.Reputacion{},
		nextID: 1,
	}
}

func (f *ReputacionFake) ListarReputacion() []models.Reputacion {
	lista := make([]models.Reputacion, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *ReputacionFake) BuscarReputacionPorID(id int) (models.Reputacion, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *ReputacionFake) CrearReputacion(r models.Reputacion) models.Reputacion {
	r.ID = f.nextID
	f.nextID++
	f.porID[r.ID] = r
	return r
}

func (f *ReputacionFake) ActualizarReputacion(id int, datos models.Reputacion) (models.Reputacion, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Reputacion{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *ReputacionFake) BorrarReputacion(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.ReputacionRepository = (*ReputacionFake)(nil)

func construirEntornoReputacion(t *testing.T) http.Handler {
	t.Helper()

	repFake := nuevoReputacionFake()
	repService := rlc.NewReputacionService(repFake)

	srv := handlers.NewServer(handlers.Deps{Reputacion: repService})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/reputacion", srv.ListarReputacion)
		r.Post("/reputacion", srv.CrearReputacion)
		r.Get("/reputacion/{id}", srv.ObtenerReputacion)
		r.Put("/reputacion/{id}", srv.ActualizarReputacion)
		r.Delete("/reputacion/{id}", srv.BorrarReputacion)
	})

	return r
}

func TestCrearReputacion_Exitoso(t *testing.T) {
	h := construirEntornoReputacion(t)
	body := `{"puntos_totales":150,"nivel":2,"acuerdo_compl":5,"calificacion_pro":4.5,"usuario_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reputacion", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCrearReputacion_UsuarioIDVacio(t *testing.T) {
	h := construirEntornoReputacion(t)
	body := `{"puntos_totales":150,"nivel":2,"acuerdo_compl":5,"calificacion_pro":4.5}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reputacion", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListarReputacion_Vacio(t *testing.T) {
	h := construirEntornoReputacion(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reputacion", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestObtenerReputacion_NoExiste(t *testing.T) {
	h := construirEntornoReputacion(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reputacion/999", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCrearReputacion_PuntosTotalesVacio(t *testing.T) {
	h := construirEntornoReputacion(t)
	body := `{"nivel":2,"acuerdo_compl":5,"calificacion_pro":4.5,"usuario_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reputacion", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
