package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"proyecto-semestral/internal/models"
)

// ListarInventario atiende GET /api/v1/inventarios.
func (s *Server) ListarInventario(w http.ResponseWriter, _ *http.Request) {
	inventarios := s.Inventario.ListarInventario() // cambiar a que todos tengan su servicio con su handler
	RespondJSON(w, http.StatusOK, inventarios)
}

// ObtenerInventario atiende GET /api/v1/inventarios/{id}.
func (s *Server) ObtenerInventario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	inventario, err := s.Inventario.BuscarInventario(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
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

	creada, err := s.Inventario.CrearInventario(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarInventario atiende PUT /api/v1/inventarios/{id}.
func (s *Server) ActualizarInventario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
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

	actualizada, err := s.Inventario.ActualizarInventario(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarInventario atiende DELETE /api/v1/inventarios/{id}.
func (s *Server) BorrarInventario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Inventario.BorrarInventario(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
