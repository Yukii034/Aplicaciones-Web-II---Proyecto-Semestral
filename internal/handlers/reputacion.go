package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarReputacion atiende GET /api/v1/reputaciones.
func (s *Server) ListarReputacion(w http.ResponseWriter, _ *http.Request) {
	reputaciones := s.Storage.ListarReputacion()
	RespondJSON(w, http.StatusOK, reputaciones)
}

// ObtenerReputacion atiende GET /api/v1/reputaciones/{id}.
func (s *Server) ObtenerReputacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	reputacion, encontrado := s.Storage.BuscarReputacionPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Inventario no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, reputacion)
}

// CrearReputacion atiende POST /api/v1/reputaciones.
func (s *Server) CrearReputacion(w http.ResponseWriter, r *http.Request) {
	var nueva models.Reputacion
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if nueva.UsuarioID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	creada := s.Storage.CrearReputacion(nueva)
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarReputacion atiende PUT /api/v1/reputaciones/{id}.
func (s *Server) ActualizarReputacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Reputacion
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if datos.UsuarioID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	actualizada, encontrada := s.Storage.ActualizarReputacion(id, datos)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "Reputacion no encontrado")
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarReputacion atiende DELETE /api/v1/reputaciones/{id}.
func (s *Server) BorrarReputacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if !s.Storage.BorrarReputacion(id) {
		RespondError(w, http.StatusNotFound, "Reputacion no encontrada")
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
