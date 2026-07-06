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

// construirEntorno arma el MISMO router que main.go (mismas rutas, mismo
// middleware.Auth real) pero con almacen en memoria y repos fake.
// Devuelve el handler listo para httptest y un token valido ya emitido.
func construirEntorno(t *testing.T) (http.Handler, string) {
	t.Helper()

	// fakes en memoria
	invFake := nuevoInventarioFake()
	usuFake := nuevoUsuarioFake()

	// servicios con los fake a agregar
	invService := pi.NewInventarioService(invFake)
	authService := service.NuevoAuthService(usuFake)

	// servidores con nil para los servicios de las demás entidades (solo para los que se están usando)
	srv := handlers.NewServer(handlers.Deps{Inventario: invService, Auth: authService})

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
			r.Get("/inventario/{id}", srv.ObtenerInventario)
			r.Put("/inventario/{id}", srv.ActualizarInventario)
			r.Delete("/inventario/{id}", srv.BorrarInventario)
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
	h.ServeHTTP(httptest.NewRecorder(), reqReg) // captura el request contra el router, no guarda la respuesta, solo que se registra

	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(cred)) // se loggea
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, reqLogin)                // guarda la respuesta ya que se necesita leer el token
	require.Equal(t, http.StatusOK, rec.Code) // recoge el codigo de estatus que mandó el servidor

	var resp struct {
		Token string `json:"token"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp)) // recoge el token en json y lo lee como string
	require.NotEmpty(t, resp.Token)                             // revisa que no esté vacío
	return resp.Toke                                            // devuelve el token como string
}

// TestCrearProducto_Exitoso: POST con token y cuerpo valido -> 201 Creado
func TestCrearInventario_Exitoso(t *testing.T) {
	h, token := construirEntorno(t) // construye el server con el token ya guardado
	body := `{"nombre":"Laptop Dell","categoria":"Tecnología","cantidad":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventario", strings.NewReader(body)) // hace un post creando un inventario nuevo
	// agrega el token al header del request, escoge la autorizacion para un bearer token, mandando luego el token como string
	// el middleware jwt espera un header con Bearer y el token
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder() // captura todo lo que el handler responda, codigo, headers y body
	h.ServeHTTP(rec, req)         // envía el request

	require.Equal(t, http.StatusCreated, rec.Code) // asegura que salga el codigo 201 de creado
}

// TestRutaProtegida_SinToken: sin header Authorization, el middleware corta
// antes de llegar al handler -> 401 Sin autorizar
func TestRutaProtegida_SinToken(t *testing.T) {
	h, _ := construirEntorno(t) // construye el server sin token gardado

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventario", nil) // hace un get para listar inventario
	// sin token a propósito
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code) // se espera que salga el codigo 401 sin autorizar, ya que no tiene token
}

func TestObtenerInventario_NoExiste(t *testing.T) {
	h, token := construirEntorno(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/inventario/999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCrearInventario_TituloVacio(t *testing.T) {
	h, token := construirEntorno(t)
	body := `{"nombre":"","categoria":"Tecnología","cantidad":1}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/inventario", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
