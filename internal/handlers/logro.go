package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarLogro atiende GET /api/v1/logros.
func (s *Server) ListarLogro(w http.ResponseWriter, _ *http.Request) {
	logros := s.Storage.ListarLogro()
	RespondJSON(w, http.StatusOK, logros)
}

// ObtenerLogro atiende GET /api/v1/logros/{id}.
func (s *Server) ObtenerLogro(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	logro, encontrado := s.Storage.BuscarLogroPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Logro no encontrada")
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
	if strings.TrimSpace(nueva.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "El campo nombre es obligatorio")
		return
	}

	creada := s.Storage.CrearLogro(nueva)
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
	if strings.TrimSpace(datos.Nombre) == "" {
		RespondError(w, http.StatusBadRequest, "El campo nombre es obligatorio")
		return
	}

	actualizada, encontrada := s.Storage.ActualizarLogro(id, datos)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "logro no encontrado")
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

	if !s.Storage.BorrarLogro(id) {
		RespondError(w, http.StatusNotFound, "Logro no encontrada")
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
