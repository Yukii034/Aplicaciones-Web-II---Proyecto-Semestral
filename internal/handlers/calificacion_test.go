package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"proyecto-semestral/internal/middleware"
	"proyecto-semestral/internal/service"

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

type acuerdoFake struct {
	porID  map[int]models.Acuerdo
	nextID int
}

func nuevoAcuerdoFake() *acuerdoFake {
	return &acuerdoFake{
		porID:  map[int]models.Acuerdo{},
		nextID: 1,
	}
}

func (f *acuerdoFake) ListarAcuerdos() []models.Acuerdo {
	lista := make([]models.Acuerdo, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *acuerdoFake) BuscarAcuerdoPorID(id int) (models.Acuerdo, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *acuerdoFake) CrearAcuerdo(a models.Acuerdo) models.Acuerdo {
	a.ID = f.nextID
	f.nextID++
	f.porID[a.ID] = a
	return a
}

func (f *acuerdoFake) ActualizarAcuerdo(id int, datos models.Acuerdo) (models.Acuerdo, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Acuerdo{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *acuerdoFake) BorrarAcuerdo(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.AcuerdoRepository = (*acuerdoFake)(nil)

func construirEntornoCalificacion(t *testing.T) (http.Handler, string) {
	t.Helper()

	caFake := nuevoCalificacionFake()
	usuFake := nuevoUsuarioFake() // viene de inventario_test.go
	acuFake := nuevoAcuerdoFake() // para la relacion

	caService := rlc.NewCalificacionService(caFake, usuFake, acuFake)
	authService := service.NuevoAuthService(usuFake)

	srv := handlers.NewServer(handlers.Deps{
		Calificacion: caService,
		Auth:         authService,
	})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Get("/calificaciones", srv.ListarCalificacion)
			r.Post("/calificaciones", srv.CrearCalificacion)
			r.Get("/calificaciones/{id}", srv.ObtenerCalificacion)
			r.Put("/calificaciones/{id}", srv.ActualizarCalificacion)
			r.Delete("/calificaciones/{id}", srv.BorrarCalificacion)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

func TestCrearCalificacion_Exitoso(t *testing.T) {
	caFake := nuevoCalificacionFake()
	usuFake := nuevoUsuarioFake()
	acuFake := nuevoAcuerdoFake()

	caService := rlc.NewCalificacionService(caFake, usuFake, acuFake)
	authService := service.NuevoAuthService(usuFake)
	srv := handlers.NewServer(handlers.Deps{Calificacion: caService, Auth: authService})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Post("/calificaciones", srv.CrearCalificacion)
		})
	})

	// registrar crea el usuario con ID 1 automáticamente
	token := registrarYObtenerToken(t, r)

	// ahora pre-crear el acuerdo
	acuFake.CrearAcuerdo(models.Acuerdo{PublicacionID: 1, IDOfertante: 1, IDSolicitante: 1, Tipo: "intercambio"})

	body := `{"comentarios":"Excelente intercambio","usuario_id":1,"acuerdo_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calificaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCrearCalificacion_ComentarioVacio(t *testing.T) {
	h, token := construirEntornoCalificacion(t)
	body := `{"comentarios":"","usuario_id":1,"acuerdo_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calificaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListarCalificacion_Vacio(t *testing.T) {
	h, token := construirEntornoCalificacion(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calificaciones", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestObtenerCalificacion_NoExiste(t *testing.T) {
	h, token := construirEntornoCalificacion(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calificaciones/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCrearCalificacion_SinUsuarioID(t *testing.T) {
	h, token := construirEntornoCalificacion(t)
	body := `{"comentarios":"Buen intercambio","acuerdo_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calificaciones", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}
