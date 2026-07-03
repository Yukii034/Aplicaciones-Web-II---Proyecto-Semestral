package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/handlers"
	"proyecto-semestral/internal/middleware"
	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	pi "proyecto-semestral/internal/service/modulo_pi"
	"proyecto-semestral/internal/storage"
)

// Fake de publicacion
type publicacionFake struct {
	items  []models.Publicacion
	nextID int
}

func (f *publicacionFake) ListarPublicacion() []models.Publicacion { return f.items }

func (f *publicacionFake) BuscarPublicacionPorID(id int) (models.Publicacion, bool) {
	for _, item := range f.items {
		if item.ID == id {
			return item, true
		}
	}
	return models.Publicacion{}, false
}

func (f *publicacionFake) CrearPublicacion(p models.Publicacion) models.Publicacion {
	p.ID = f.nextID
	f.nextID++
	f.items = append(f.items, p)
	return p
}

func (f *publicacionFake) ActualizarPublicacion(id int, datos models.Publicacion) (models.Publicacion, bool) {
	for idx, item := range f.items {
		if item.ID == id {
			datos.ID = id
			f.items[idx] = datos
			return datos, true
		}
	}
	return models.Publicacion{}, false
}

func (f *publicacionFake) BorrarPublicacion(id int) bool {
	for idx, item := range f.items {
		if item.ID == id {
			f.items = append(f.items[:idx], f.items[idx+1:]...)
			return true
		}
	}
	return false
}

var _ storage.PublicacionRepository = (*publicacionFake)(nil)

// construirEntornoP arma el router solo con publicacion y auth
func construirEntornoP(t *testing.T) (http.Handler, string) {
	t.Helper()

	pubFake := &publicacionFake{nextID: 1}
	usuFake := nuevoUsuarioFake() // viene de inventario_test.go

	pubService := pi.NewPublicacionService(pubFake)
	authService := service.NuevoAuthService(usuFake)

	srv := handlers.NewServer(handlers.Deps{
		Publicacion: pubService,
		Auth:        authService,
	})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Get("/publicaciones", srv.ListarPublicacion)
			r.Post("/publicaciones", srv.CrearPublicacion)
			r.Get("/publicaciones/{id}", srv.ObtenerPublicacion)
			r.Put("/publicaciones/{id}", srv.ActualizarPublicacion)
			r.Delete("/publicaciones/{id}", srv.BorrarPublicacion)
		})
	})

	token := registrarYObtenerToken(t, r) // viene de inventario_test.go
	return r, token
}

// TestCrearPublicacion_Exitoso: POST con token y cuerpo valido -> 201 Creado
func TestCrearPublicacion_Exitoso(t *testing.T) {
	h, token := construirEntornoP(t)
	body := `{"titulo":"Cambio laptop por tablet","tipo_oferta":"intercambio","estado_publicacion":"disponible","inventario_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/publicaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

// TestRutaProtegida_SinTokenP: sin token -> 401
func TestRutaProtegida_SinTokenP(t *testing.T) {
	h, _ := construirEntornoP(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/publicaciones", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestObtenerPublicacion_NoExiste(t *testing.T) {
	h, token := construirEntornoP(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/publicaciones/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCrearPublicacion_TituloVacio(t *testing.T) {
	h, token := construirEntornoP(t)
	body := `{"titulo":"","tipo_oferta":"intercambio"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/publicaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBorrarPublicacion_NoExiste(t *testing.T) {
	h, token := construirEntornoP(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/publicaciones/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestActualizarPublicacion_NoExiste(t *testing.T) {
	h, token := construirEntornoP(t)
	body := `{"titulo":"Cambio tablet","tipo_oferta":"intercambio"}`

	req := httptest.NewRequest(http.MethodPut, "/api/v1/publicaciones/999", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
