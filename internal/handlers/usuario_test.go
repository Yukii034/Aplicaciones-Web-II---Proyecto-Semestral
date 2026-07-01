package handlers_test

import (
	"encoding/json"
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
	aiu "proyecto-semestral/internal/service/modulo_aiu" // <-- Ahora sí con el alias aiu
	"proyecto-semestral/internal/storage"
)

// --- Fake de Usuario (Sincronizado por Email e ID) ---
type usuarioRepoFake struct {
	porEmail map[string]models.Usuario
	porID    map[int]models.Usuario
	nextID   int
}

func nuevoUsuarioRepoFake() *usuarioRepoFake {
	return &usuarioRepoFake{
		porEmail: map[string]models.Usuario{},
		porID:    map[int]models.Usuario{},
		nextID:   1,
	}
}

func (f *usuarioRepoFake) ListarUsuarios() []models.Usuario {
	lista := make([]models.Usuario, 0, len(f.porID))
	for _, u := range f.porID {
		lista = append(lista, u)
	}
	return lista
}

func (f *usuarioRepoFake) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	u, ok := f.porID[id]
	return u, ok
}

func (f *usuarioRepoFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

func (f *usuarioRepoFake) CrearUsuario(usuario models.Usuario) models.Usuario {
	usuario.ID = f.nextID
	f.nextID++
	f.porEmail[usuario.Email] = usuario
	f.porID[usuario.ID] = usuario
	return usuario
}

func (f *usuarioRepoFake) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Usuario{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	f.porEmail[datos.Email] = datos
	return datos, true
}

func (f *usuarioRepoFake) BorrarUsuario(id int) bool {
	u, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	delete(f.porEmail, u.Email)
	return true
}

var _ storage.UserRepository = (*usuarioRepoFake)(nil)

// --- Constructor del Entorno de Pruebas ---
func construirEntornoUsuarios(t *testing.T) (http.Handler, string) {
	t.Helper()

	// 1. Instanciar el repositorio fake en memoria
	usuFake := nuevoUsuarioFake()

	// 2. Instanciar los servicios del sistema
	authService := service.NuevoAuthService(usuFake)

	// Usando el alias 'aiu' correctamente
	uService := aiu.NewUsuarioService(usuFake)

	// 3. Crear el servidor pasando tu servicio AIU
	srv := handlers.NewServer(handlers.Deps{Usuario: uService, Auth: authService})

	// 4. Configurar el Router con Middleware Real y tus endpoints
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			// Endpoints de tu Módulo AIU
			r.Get("/usuarios", srv.ListarUsuarios)
			r.Get("/usuarios/{id}", srv.ObtenerUsuario)
			r.Post("/usuarios", srv.CrearUsuario)
			r.Put("/usuarios/{id}", srv.ActualizarUsuario)
			r.Delete("/usuarios/{id}", srv.EliminarUsuario)
		})
	})

	token := registrarYObtenerTokenUsuarios(t, r)
	return r, token
}

// Helper para automatizar el registro, login y extracción del Token JWT real
func registrarYObtenerTokenUsuarios(t *testing.T, h http.Handler) string {
	t.Helper()
	cred := `{"email":"test@test.com","password":"secreta123"}`

	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(cred))
	h.ServeHTTP(httptest.NewRecorder(), reqReg)

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(cred))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqLogin)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.NotEmpty(t, resp.Token)
	return resp.Token
}

// --- Casos de Prueba ---

func TestCrearUsuario_Exitoso(t *testing.T) {
	h, token := construirEntornoUsuarios(t)
	body := `{"nombre":"Juan Perez","email":"juan@perez.com","tipo":"admin","ciudad":"Santiago"}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/usuarios", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var usuario models.Usuario
	err := json.NewDecoder(rec.Body).Decode(&usuario)
	require.NoError(t, err)
	assert.Equal(t, "Juan Perez", usuario.Nombre)
}

func TestListarUsuarios_Exitoso(t *testing.T) {
	h, token := construirEntornoUsuarios(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRutaUsuarios_SinToken(t *testing.T) {
	h, _ := construirEntornoUsuarios(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/usuarios", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
