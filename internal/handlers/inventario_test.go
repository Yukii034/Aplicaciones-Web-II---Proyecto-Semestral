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

// Fake de usuario (guarda de verdad para que register y login funcionen)
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

// Funciones requeridas para el fake
func (f *usuarioFake) BuscarUsuarioPorEmail(email string) (models.Usuario, bool) {
	u, ok := f.porEmail[email]
	return u, ok
}

func (f *usuarioFake) CrearUsuario(u models.Usuario) models.Usuario {
	u.ID = f.nextID
	f.nextID++
	f.porEmail[u.Email] = u
	return u
}

// las demás funciones solo existen para cumplir con la interfaz
func (f *usuarioFake) ListarUsuarios() []models.Usuario {
	return nil
}

func (f *usuarioFake) BuscarUsuarioPorID(id int) (models.Usuario, bool) {
	return models.Usuario{}, false
}

func (f *usuarioFake) ActualizarUsuario(id int, datos models.Usuario) (models.Usuario, bool) {
	return models.Usuario{}, false
}

func (f *usuarioFake) BorrarUsuario(id int) bool {
	return false
}

// Verifica que usuarioFake cumple todos los métodos de UserRepository, es para saber si compila correctamente
var _ storage.UserRepository = (*usuarioFake)(nil)

// Fake de inventario
type inventarioFake struct {
	items  []models.Inventario
	nextID int
}

func nuevoInventarioFake() *inventarioFake {
	return &inventarioFake{nextID: 1}
}

func (f *inventarioFake) ListarInventario() []models.Inventario {
	return f.items
}

func (f *inventarioFake) BuscarInventarioPorID(id int) (models.Inventario, bool) {
	for _, item := range f.items {
		if item.ID == id {
			return item, true
		}
	}
	return models.Inventario{}, false
}

func (f *inventarioFake) CrearInventario(i models.Inventario) models.Inventario {
	i.ID = f.nextID
	f.nextID++
	f.items = append(f.items, i)
	return i
}

func (f *inventarioFake) ActualizarInventario(id int, datos models.Inventario) (models.Inventario, bool) {
	return models.Inventario{}, false
}

func (f *inventarioFake) BorrarInventario(id int) bool {
	return false
}

var _ storage.InventarioRepository = (*inventarioFake)(nil)

// Fake de publicacion
type publicacionFake struct {
	items  []models.Publicacion
	nextID int
}

func nuevoPublicacionFake() *publicacionFake {
	return &publicacionFake{nextID: 1}
}

func (f *publicacionFake) ListarPublicacion() []models.Publicacion {
	return f.items
}

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
	return models.Publicacion{}, false
}

func (f *publicacionFake) BorrarPublicacion(id int) bool {
	return false
}

var _ storage.PublicacionRepository = (*publicacionFake)(nil)

// construirEntorno arma el MISMO router que main.go (mismas rutas, mismo
// middleware.Auth real) pero con almacen en memoria y repos fake.
// Devuelve el handler listo para httptest y un token valido ya emitido.
func construirEntorno(t *testing.T) (http.Handler, string) {
	t.Helper()

	// fakes en memoria
	invFake := nuevoInventarioFake()
	pubFake := nuevoPublicacionFake()
	usuFake := nuevoUsuarioFake()

	// servicios con los fake a agregar
	invService := pi.NewInventarioService(invFake)
	pubService := pi.NewPublicacionService(pubFake)
	authService := service.NewAuthService(usuFake)

	// servidores con nil para los servicios de las demás entidades (solo para los que se están usando)
	srv := handlers.NewServer(invService, pubService, nil, nil, nil, nil, nil, nil, nil, authService)

	// router completo con middleware real
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", srv.Registrar)
		r.Post("/auth/login", srv.Login)

		// rutas protegidas con autenticación
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))
			r.Get("/inventario", srv.ListarInventario)
			r.Post("/inventario", srv.CrearInventario)
		})
	})

	token := registrarYObtenerToken(t, r)
	return r, token
}

// hace el proceso de registro de usuario y se obtiene el token requerido para acceder a rutas protegidas
func registrarYObtenerToken(t *testing.T, h http.Handler) string {
	t.Helper() // para que go sepa que es una helper de un test, este apunta a quién llamó a esta funcion y no a la propia funcion
	cred := `{"email":"test@test.com","password":"secreta123"}`

	// registra simulando un post y convierte el json en un lector que el request puede enviar como body
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(cred))
	h.ServeHTTP(httptest.NewRecorder(), reqReg) // captura el request contra el router

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(cred)) // se loggea
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqLogin)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp)) // recoge el toquen en json y lo cambia a string
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
