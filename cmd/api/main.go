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
	if err := db.AutoMigrate(&models.Inventario{}, &models.Publicacion{}); err != nil {
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
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
