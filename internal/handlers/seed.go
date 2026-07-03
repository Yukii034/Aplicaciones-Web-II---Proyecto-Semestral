package handlers

import "net/http"

// Sembrar inserta datos de ejemplo si la base esta vacia. Pensado para
// dispararse manualmente una vez despues de levantar el contenedor.
func (s *Server) Sembrar(w http.ResponseWriter, r *http.Request) {
	almacenSQLite, ok := s.Almacen.(interface{ SembrarSiVacio() })
	if !ok {
		http.Error(w, "seed no disponible con este backend", http.StatusNotImplemented)
		return
	}
	almacenSQLite.SembrarSiVacio()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Seed ejecutado (o ya existian datos)"))
}
