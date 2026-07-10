package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"proyecto-semestral/internal/handlers"
	"proyecto-semestral/internal/middleware"
	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/service"
	rlc "proyecto-semestral/internal/service/modulo_rlc"
	"proyecto-semestral/internal/storage"
)

type Logro_UsuarioFake struct {
	porID  map[int]models.Logro_Usuario
	nextID int
}

func nuevoLogro_UsuarioFake() *Logro_UsuarioFake {
	return &Logro_UsuarioFake{
		porID:  map[int]models.Logro_Usuario{},
		nextID: 1,
	}
}

func (f *Logro_UsuarioFake) ListarLogro_Usuario() []models.Logro_Usuario {
	lista := make([]models.Logro_Usuario, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *Logro_UsuarioFake) BuscarLogro_UsuarioPorID(id int) (models.Logro_Usuario, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *Logro_UsuarioFake) CrearLogro_Usuario(lu models.Logro_Usuario) models.Logro_Usuario {
	lu.ID = f.nextID
	f.nextID++
	f.porID[lu.ID] = lu
	return lu
}

func (f *Logro_UsuarioFake) ActualizarLogro_Usuario(id int, datos models.Logro_Usuario) (models.Logro_Usuario, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Logro_Usuario{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *Logro_UsuarioFake) BorrarLogro_Usuario(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.Logro_UsuarioRepository = (*Logro_UsuarioFake)(nil)

func construirEntornoLogro_Usuario(t *testing.T) (http.Handler, string) {
	t.Helper()

	luFake := nuevoLogro_UsuarioFake()
	usuFake := nuevoUsuarioFake() // viene de inventario_test.go
	logFake := nuevoLogroFake()   // viene de logro_test.go

	luService := rlc.NewLogro_UsuarioService(luFake, logFake, usuFake)

	authService := service.NuevoAuthService(usuFake)

	srv := handlers.NewServer(handlers.Deps{
		Logro_Usuario: luService,
		Auth:          authService,
	})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Get("/logro_usuarios", srv.ListarLogro_Usuario)
			r.Post("/logro_usuarios", srv.CrearLogro_Usuario)
			r.Get("/logro_usuarios/{id}", srv.ObtenerLogro_Usuario)
			r.Put("/logro_usuarios/{id}", srv.ActualizarLogro_Usuario)
			r.Delete("/logro_usuarios/{id}", srv.BorrarLogro_Usuario)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

func TestCrearLogro_Usuario_Exitoso(t *testing.T) {
	luFake := nuevoLogro_UsuarioFake()
	usuFake := nuevoUsuarioFake()
	logFake := nuevoLogroFake()

	luService := rlc.NewLogro_UsuarioService(luFake, logFake, usuFake)
	authService := service.NuevoAuthService(usuFake)
	srv := handlers.NewServer(handlers.Deps{Logro_Usuario: luService, Auth: authService})

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Post("/logro_usuarios", srv.CrearLogro_Usuario)
		})
	})

	// registrar crea el usuario con ID 1 automáticamente
	token := registrarYObtenerToken(t, r)

	// ahora pre-crear el logro
	logFake.CrearLogro(models.Logro{Nombre: "Primer Intercambio", PuntosRequeridos: 50})

	body := `{"fechas_desbloqueado":"2026-01-01T00:00:00Z","usuario_id":1,"logro_id":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/logro_usuarios", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

func TestCrearLogro_Usuario_FechaVacia(t *testing.T) {
	h, token := construirEntornoLogro_Usuario(t)
	body := `{"usuario_id":1,"logro_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logro_usuarios", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListarLogro_Usuario_Vacio(t *testing.T) {
	h, token := construirEntornoLogro_Usuario(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logro_usuarios", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

func TestObtenerLogro_Usuario_NoExiste(t *testing.T) {
	h, token := construirEntornoLogro_Usuario(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/logro_usuarios/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCrearLogro_Usuario_FechaVaciaExplicita(t *testing.T) {
	h, token := construirEntornoLogro_Usuario(t)
	body := `{"fechas_desbloqueado":"0001-01-01T00:00:00Z","usuario_id":1,"logro_id":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/logro_usuarios", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

// solo para que el compilador no se queje de time importado sin usar
var _ = time.Now
