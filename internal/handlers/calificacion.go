package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarCalificacion atiende GET /api/v1/calificaciones.
func (s *Server) ListarCalificacion(w http.ResponseWriter, _ *http.Request) {
	calificaciones := s.Storage.ListarCalificacion()
	RespondJSON(w, http.StatusOK, calificaciones)
}

// ObtenerCalificacion atiende GET /api/v1/calificaciones/{id}.
func (s *Server) ObtenerCalificacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	Calificacion, encontrado := s.Storage.BuscarCalificacionPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Calificacion no encontrada")
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

	if strings.TrimSpace(nueva.Comentarios) == "" {
		RespondError(w, http.StatusBadRequest, "El campo comentario es obligatorio")
		return
	}

	if nueva.UsuarioID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	creada := s.Storage.CrearCalificacion(nueva)
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarCalificacion atiende PUT /api/v1/calificaciones/{id}.
func (s *Server) ActualizarCalificacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Calificacion
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	if strings.TrimSpace(datos.Comentarios) == "" {
		RespondError(w, http.StatusBadRequest, "El campo comentarios es obligatorio")
		return
	}

	if datos.UsuarioID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	actualizada, encontrada := s.Storage.ActualizarCalificacion(id, datos)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "Calificacion no encontrado")
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarLogro atiende DELETE /api/v1/logros/{id}.
func (s *Server) BorrarCalificacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if !s.Storage.BorrarCalificacion(id) {
		RespondError(w, http.StatusNotFound, "Calificacion no encontrada")
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
