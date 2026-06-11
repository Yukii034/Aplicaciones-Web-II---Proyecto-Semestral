package main

import (
	"log"
	"net/http"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"proyecto-semestral/internal/handlers"
	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/storage"
)

func main() {
	// 1. Abrir SQLite y migrar
	db, err := gorm.Open(sqlite.Open("proyecto.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo abrir la base de datos: ", err)
	}
	if err := db.AutoMigrate(
		&models.Inventario{},
		&models.Publicacion{},
		&models.Reputacion{},
		&models.Logro_Usuario{},
		&models.Logro{},
		&models.Calificacion{},
		&models.Usuario{},
		&models.Acuerdo{},
		&models.AcuerdoItem{},
	); err != nil {
		log.Fatal("falló AutoMigrate: ", err)
	}

	// 2. Crear almacén y sembrar datos de ejemplo
	almacen := storage.NuevoAlmacenSQLite(db)
	almacen.SembrarSiVacio()

	// 3. Crear servidor
	servidor := handlers.NewServer(almacen)

	// 4. Router
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// 5. Rutas
	r.Route("/api/v1", func(r chi.Router) {
		// Módulo de Publicaciones e Inventario
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

		// Módulo de Reputación y Logros
		r.Get("/reputaciones", servidor.ListarReputacion)
		r.Post("/reputaciones", servidor.CrearReputacion)
		r.Get("/reputaciones/{id}", servidor.ObtenerReputacion)
		r.Put("/reputaciones/{id}", servidor.ActualizarReputacion)
		r.Delete("/reputaciones/{id}", servidor.BorrarReputacion)

		r.Get("/logro_usuarios", servidor.ListarLogro_Usuario)
		r.Post("/logro_usuarios", servidor.CrearLogro_Usuario)
		r.Get("/logro_usuarios/{id}", servidor.ObtenerLogro_Usuario)
		r.Put("/logro_usuarios/{id}", servidor.ActualizarLogro_Usuario)
		r.Delete("/logro_usuarios/{id}", servidor.BorrarLogro_Usuario)

		r.Get("/logros", servidor.ListarLogro)
		r.Post("/logros", servidor.CrearLogro)
		r.Get("/logros/{id}", servidor.ObtenerLogro)
		r.Put("/logros/{id}", servidor.ActualizarLogro)
		r.Delete("/logros/{id}", servidor.BorrarLogro)

		r.Get("/calificaciones", servidor.ListarCalificacion)
		r.Post("/calificaciones", servidor.CrearCalificacion)
		r.Get("/calificaciones/{id}", servidor.ObtenerCalificacion)
		r.Put("/calificaciones/{id}", servidor.ActualizarCalificacion)
		r.Delete("/calificaciones/{id}", servidor.BorrarCalificacion)

		// Módulo de Acuerdos y Transacciones
		r.Get("/usuarios", servidor.ListarUsuarios)
		r.Post("/usuarios", servidor.CrearUsuario)
		r.Get("/usuarios/{id}", servidor.ObtenerUsuario)
		r.Put("/usuarios/{id}", servidor.ActualizarUsuario)
		r.Delete("/usuarios/{id}", servidor.EliminarUsuario)

		r.Get("/acuerdos", servidor.ListarAcuerdo)
		r.Post("/acuerdos", servidor.CrearAcuerdo)
		r.Get("/acuerdos/{id}", servidor.ObtenerAcuerdo)
		r.Put("/acuerdos/{id}", servidor.ActualizarAcuerdo)
		r.Delete("/acuerdos/{id}", servidor.EliminarAcuerdo)

		r.Get("/acuerdo_items", servidor.ListarAcuerdoItem)
		r.Post("/acuerdo_items", servidor.CrearAcuerdoItem)
		r.Get("/acuerdo_items/{id}", servidor.ObtenerAcuerdoItem)
		r.Put("/acuerdo_items/{id}", servidor.ActualizarAcuerdoItem)
		r.Delete("/acuerdo_items/{id}", servidor.EliminarAcuerdoItem)
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
