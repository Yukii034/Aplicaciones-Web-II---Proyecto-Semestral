package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"proyecto-semestral/internal/config"
	"proyecto-semestral/internal/handlers"
	"proyecto-semestral/internal/httpserver"
	"proyecto-semestral/internal/middleware"
	"proyecto-semestral/internal/service"
	aiu "proyecto-semestral/internal/service/modulo_aiu"
	pi "proyecto-semestral/internal/service/modulo_pi"
	rlc "proyecto-semestral/internal/service/modulo_rlc"
	"proyecto-semestral/internal/storage"
)

func main() {
	cfg := config.Cargar()
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config.Config) error {
	// 1. Recursos de almacenamiento: abre DB (sqlite local o postgres en Docker), migra.
	recursos, err := storage.Inicializar(cfg.DBDriver, cfg.DBDsn, cfg.RutaDB)
	if err != nil {
		return err
	}
	defer func() { _ = recursos.Cerrar() }()
	log.Printf("Motor de base de datos: %s | Backend: %s", cfg.DBDriver, recursos.BackendUsado)

	// 2. Capa de servicio.
	authService := service.NuevoAuthService(recursos.Usuarios)
	invService := pi.NewInventarioService(recursos.Almacen)
	pubService := pi.NewPublicacionService(recursos.Almacen, recursos.Almacen, recursos.Usuarios)
	repService := rlc.NewReputacionService(recursos.Almacen)
	luService := rlc.NewLogro_UsuarioService(recursos.Almacen)
	lService := rlc.NewLogroService(recursos.Almacen)
	caService := rlc.NewCalificacionService(recursos.Almacen)
	acIService := aiu.NewAcuerdoService(recursos.Almacen)
	acService := aiu.NewAcuerdoItemService(recursos.Almacen)
	userService := aiu.NewUsuarioService(recursos.Almacen)

	authMW := middleware.Auth(authService)

	servidor := handlers.NewServer(handlers.Deps{
		Almacen:       recursos.Almacen,
		Inventario:    invService,
		Publicacion:   pubService,
		Reputacion:    repService,
		Logro_Usuario: luService,
		Logro:         lService,
		Calificacion:  caService,
		Acuerdo:       acIService,
		AcuerdoItem:   acService,
		Usuario:       userService,
		Auth:          authService,
	})

	// 3. Router + middleware global.
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.CORS)

	// 4. Rutas.
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", servidor.Registrar)
		r.Post("/auth/login", servidor.Login)
		r.Post("/seed", servidor.Sembrar) // para sembrar datos si está vacío

		r.Group(func(r chi.Router) {
			r.Use(authMW) // autenticacion para usuario registrado

			r.Get("/inventario", servidor.ListarInventario)
			r.Post("/inventario", servidor.CrearInventario)
			r.Get("/inventario/{id}", servidor.ObtenerInventario)
			r.Put("/inventario/{id}", servidor.ActualizarInventario)
			r.Delete("/inventario/{id}", servidor.BorrarInventario)

			r.Get("/publicaciones", servidor.ListarPublicacion)
			r.Post("/publicaciones", servidor.CrearPublicacion)
			r.Get("/publicaciones/{id}", servidor.ObtenerPublicacion)
			r.Put("/publicaciones/{id}", servidor.ActualizarPublicacion)
			r.Delete("/publicaciones/{id}", servidor.BorrarPublicacion)

			r.Get("/reputaciones", servidor.ListarReputacion)
			r.Post("/reputaciones", servidor.CrearReputacion)
			r.Get("/reputaciones/{id}", servidor.ObtenerReputacion)

			r.Group(func(r chi.Router) {
				r.Use(middleware.SoloAdmin) // autenticacion de admin
				r.Put("/reputaciones/{id}", servidor.ActualizarReputacion)
				r.Delete("/reputaciones/{id}", servidor.BorrarReputacion)
			})

			r.Get("/logro_usuarios", servidor.ListarLogro_Usuario)
			r.Get("/logro_usuarios/{id}", servidor.ObtenerLogro_Usuario)
			r.Group(func(r chi.Router) {
				r.Use(middleware.SoloAdmin)
				r.Post("/logro_usuarios", servidor.CrearLogro_Usuario)
				r.Put("/logro_usuarios/{id}", servidor.ActualizarLogro_Usuario)
				r.Delete("/logro_usuarios/{id}", servidor.BorrarLogro_Usuario)
			})

			r.Get("/logros", servidor.ListarLogro)
			r.Get("/logros/{id}", servidor.ObtenerLogro)
			r.Group(func(r chi.Router) {
				r.Use(middleware.SoloAdmin)
				r.Post("/logros", servidor.CrearLogro)
				r.Put("/logros/{id}", servidor.ActualizarLogro)
				r.Delete("/logros/{id}", servidor.BorrarLogro)
			})

			r.Get("/calificaciones", servidor.ListarCalificacion)
			r.Post("/calificaciones", servidor.CrearCalificacion)
			r.Get("/calificaciones/{id}", servidor.ObtenerCalificacion)
			r.Put("/calificaciones/{id}", servidor.ActualizarCalificacion)
			r.Delete("/calificaciones/{id}", servidor.BorrarCalificacion)

			r.Group(func(r chi.Router) {
				r.Use(middleware.SoloAdmin)
				r.Get("/usuarios", servidor.ListarUsuarios)
				r.Post("/usuarios", servidor.CrearUsuario)
				r.Get("/usuarios/{id}", servidor.ObtenerUsuario)
				r.Put("/usuarios/{id}", servidor.ActualizarUsuario)
				r.Delete("/usuarios/{id}", servidor.EliminarUsuario)
			})

			r.Get("/acuerdos", servidor.ListarAcuerdo)
			r.Post("/acuerdos", servidor.CrearAcuerdo)
			r.Get("/acuerdos/{id}", servidor.ObtenerAcuerdo)
			r.Put("/acuerdos/{id}", servidor.ActualizarAcuerdo) // ofertante/solicitante avanzan el flujo del trato
			r.Delete("/acuerdos/{id}", servidor.EliminarAcuerdo)

			r.Get("/acuerdo_items", servidor.ListarAcuerdoItem)
			r.Post("/acuerdo_items", servidor.CrearAcuerdoItem)
			r.Get("/acuerdo_items/{id}", servidor.ObtenerAcuerdoItem)
			r.Put("/acuerdo_items/{id}", servidor.ActualizarAcuerdoItem) // ajustar items mientras se negocia
			r.Delete("/acuerdo_items/{id}", servidor.EliminarAcuerdoItem)
		})
	})

	// 5. Servidor HTTP con graceful shutdown.
	srv := httpserver.Nuevo(
		r,
		httpserver.ConPuerto(cfg.Puerto),
		httpserver.ConReadTimeout(cfg.ReadTimeout),
		httpserver.ConWriteTimeout(cfg.WriteTimeout),
	)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errServidor := make(chan error, 1)
	go func() {
		log.Println("Servidor escuchando en http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errServidor <- err
		}
	}()

	select {
	case err := <-errServidor:
		return err
	case <-ctx.Done():
		log.Println("Senal de apagado recibida, cerrando ordenadamente...")
	}

	ctxApagado, cancelar := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelar()
	if err := srv.Shutdown(ctxApagado); err != nil {
		return err
	}
	log.Println("Servidor detenido limpiamente.")
	return nil
}
