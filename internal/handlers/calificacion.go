package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto-semestral/internal/models"
)

// ListarCalificacion atiende GET /api/v1/calificaciones.
func (s *Server) ListarCalificacion(w http.ResponseWriter, _ *http.Request) {
	calificaciones := s.Calificacion.ListarCalificacion()
	RespondJSON(w, http.StatusOK, calificaciones)
}

// ObtenerCalificacion atiende GET /api/v1/calificaciones/{id}.
func (s *Server) ObtenerCalificacion(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	Calificacion, err := s.Calificacion.BuscarCalificacion(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, Calificacion)
}

// CrearCalificacion atiende POST /api/v1/calificaciones.
func (s *Server) CrearCalificacion(w http.ResponseWriter, r *http.Request) {
	var nueva models.Calificacion
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	creada, err := s.Calificacion.CrearCalificacion(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarCalificacion atiende PUT /api/v1/calificaciones/{id}.
func (s *Server) ActualizarCalificacion(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Calificacion
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	actualizada, err := s.Calificacion.ActualizarCalificacion(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarLogro atiende DELETE /api/v1/logros/{id}.
func (s *Server) BorrarCalificacion(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Calificacion.BorrarCalificacion(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
