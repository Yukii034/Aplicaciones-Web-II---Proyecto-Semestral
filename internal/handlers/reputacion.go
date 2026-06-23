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
	reputaciones := s.Reputacion.ListarReputacion()
	RespondJSON(w, http.StatusOK, reputaciones)
}

// ObtenerReputacion atiende GET /api/v1/reputaciones/{id}.
func (s *Server) ObtenerReputacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	reputacion, err := s.Reputacion.BuscarReputacion(id)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Reputacion no encontrada")
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

	creada, err := s.Reputacion.CrearReputacion(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
	}
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

	actualizada, err := s.Reputacion.ActualizarReputacion(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
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

	if err := s.Reputacion.BorrarReputacion(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
