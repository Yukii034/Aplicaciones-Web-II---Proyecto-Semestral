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
	pi "proyecto-semestral/internal/service/modulo_pi"
	"proyecto-semestral/internal/storage"
)

// usuarioRepoFake: repositorio de usuarios en memoria para los tests de handler.
// --- Fake de usuario (guarda de verdad para que register y login funcionen) ---
type usuarioFake struct {
	porEmail map[string]models.Usuario
	nextID   int
}

func nuevoUsuarioFake() *usuarioFake {
	return &usuarioFake{
		porEmail: map[string]models.Usuario{},
		nextID:   1,
	}
}

func (f *usuarioFake) ListarUsuarios() []models.Usuario {
	return nil
}

func (f *usuarioFake) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	return models.Usuario{}, false
}

func (f *usuarioFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

func (f *usuarioFake) CrearUsuario(usuario models.Usuario) models.Usuario {
	usuario.ID = f.nextID
	f.nextID++
	f.porEmail[usuario.Email] = usuario
	return usuario
}

func (f *usuarioFake) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}

func (f *usuarioFake) BorrarUsuario(id int) bool {
	return false
}

var _ storage.UserRepository = (*usuarioFake)(nil)

type inventarioFake struct {
	porID  map[int]models.Inventario
	nextID int
}

func nuevoInventarioFake() *inventarioFake {
	return &inventarioFake{
		porID:  map[int]models.Inventario{},
		nextID: 1,
	}
}

func (f *inventarioFake) ListarInventario() []models.Inventario {
	lista := make([]models.Inventario, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *inventarioFake) BuscarInventarioPorID(id int) (models.Inventario, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *inventarioFake) CrearInventario(inventario models.Inventario) models.Inventario {
	inventario.ID = f.nextID
	f.nextID++
	f.porID[inventario.ID] = inventario
	return inventario
}

func (f *inventarioFake) ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Inventario{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *inventarioFake) BorrarInventario(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.InventarioRepository = (*inventarioFake)(nil)

type publicacionFake struct {
	porID  map[int]models.Publicacion
	nextID int
}

func nuevoPublicacionFake() *publicacionFake {
	return &publicacionFake{
		porID:  map[int]models.Publicacion{},
		nextID: 1,
	}
}

func (f *publicacionFake) ListarPublicaciones() []models.Publicacion {
	lista := make([]models.Publicacion, 0, len(f.porID))
	for _, item := range f.porID {
		lista = append(lista, item)
	}
	return lista
}

func (f *publicacionFake) ListarPublicacion() []models.Publicacion {
	return f.ListarPublicaciones()
}

func (f *publicacionFake) BuscarPublicacionPorID(id int) (models.Publicacion, bool) {
	item, ok := f.porID[id]
	return item, ok
}

func (f *publicacionFake) CrearPublicacion(publicacion models.Publicacion) models.Publicacion {
	publicacion.ID = f.nextID
	f.nextID++
	f.porID[publicacion.ID] = publicacion
	return publicacion
}

func (f *publicacionFake) ActualizarPublicacion(id int, datos models.Publicacion) (models.Publicacion, bool) {
	_, ok := f.porID[id]
	if !ok {
		return models.Publicacion{}, false
	}
	datos.ID = id
	f.porID[id] = datos
	return datos, true
}

func (f *publicacionFake) BorrarPublicacion(id int) bool {
	_, ok := f.porID[id]
	if !ok {
		return false
	}
	delete(f.porID, id)
	return true
}

var _ storage.PublicacionRepository = (*publicacionFake)(nil)

// construirEntorno arma el MISMO router que main.go (mismas rutas, mismo
// middleware.Auth real) pero con almacen en memoria y repo de usuarios fake.
// Devuelve el handler listo para httptest y un token valido ya emitido.
//
// Clave pedagogica: probamos a traves del middleware REAL, no de uno simplificado.
// Si el wiring de la ruta protegida se rompe, este test se entera.
func construirEntorno(t *testing.T) (http.Handler, string) {
	t.Helper()

	// fakes en memoria
	invFake := nuevoInventarioFake() // supuestamente esto se arregla si se arreglan las otras entidades xd
	pubFake := nuevoPublicacionFake()
	usuFake := nuevoUsuarioFake()

	// services
	invService := pi.NewInventarioService(invFake)
	pubService := pi.NewPublicacionService(pubFake)
	authService := service.NewAuthService(usuFake)

	// servidor con nil para los servicios de tus compañeros
	srv := handlers.NewServer(invService, pubService, nil, nil, nil, nil, nil, nil, nil, authService)

	// router completo con middleware real
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Get("/inventario", srv.ListarInventario)
			r.Post("/inventario", srv.CrearInventario)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

func registrarYObtenerToken(t *testing.T, h http.Handler) string {
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

// TestCrearProducto_Exitoso: POST con token y cuerpo valido -> 201 Created.
func TestCrearInventario_Exitoso(t *testing.T) {
	h, token := construirEntorno(t)
	body := `{"nombre":"Laptop Dell","categoria":"Tecnología","cantidad":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventario", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}

// TestRutaProtegida_SinToken: sin header Authorization, el middleware corta
// antes de llegar al handler -> 401 Unauthorized.
func TestRutaProtegida_SinToken(t *testing.T) {
	h, _ := construirEntorno(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventario", nil)
	// sin token a propósito
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
