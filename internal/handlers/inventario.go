package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarInventario atiende GET /api/v1/inventarios.
func (s *Server) ListarInventario(w http.ResponseWriter, _ *http.Request) {
	inventarios := s.Storage.ListarInventario()
	RespondJSON(w, http.StatusOK, inventarios)
}

// ObtenerInventario atiende GET /api/v1/inventarios/{id}.
func (s *Server) ObtenerInventario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	inventario, encontrado := s.Storage.BuscarInventarioPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Inventario no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, inventario)
}

// CrearInventario atiende POST /api/v1/inventarios.
func (s *Server) CrearInventario(w http.ResponseWriter, r *http.Request) {
	var nueva models.Inventario
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if strings.TrimSpace(nueva.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "El campo nombre es obligatorio")
		return
	}

	creada := s.Storage.CrearInventario(nueva)
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarInventario atiende PUT /api/v1/inventarios/{id}.
func (s *Server) ActualizarInventario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Inventario
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if strings.TrimSpace(datos.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "El campo título es obligatorio")
		return
	}

	actualizada, encontrada := s.Storage.ActualizarInventario(id, datos)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "Inventario no encontrado")
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarInventario atiende DELETE /api/v1/inventarios/{id}.
func (s *Server) BorrarInventario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if !s.Storage.BorrarInventario(id) {
		RespondError(w, http.StatusNotFound, "Inventario no encontrado")
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
