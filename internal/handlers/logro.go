package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarLogro atiende GET /api/v1/logros.
func (s *Server) ListarLogro(w http.ResponseWriter, _ *http.Request) {
	logros := s.Logro.ListarLogro()
	RespondJSON(w, http.StatusOK, logros)
}

// ObtenerLogro atiende GET /api/v1/logros/{id}.
func (s *Server) ObtenerLogro(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	logro, err := s.Logro.BuscarLogro(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, logro)
}

// CrearLogro atiende POST /api/v1/logros.
func (s *Server) CrearLogro(w http.ResponseWriter, r *http.Request) {
	var nueva models.Logro
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	creada, err := s.Logro.CrearLogro(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarLogro atiende PUT /api/v1/logros/{id}.
func (s *Server) ActualizarLogro(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Logro
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	actualizada, err := s.Logro.ActualizarLogro(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarLogro atiende DELETE /api/v1/logros/{id}.
func (s *Server) BorrarLogro(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Logro.BorrarLogro(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
