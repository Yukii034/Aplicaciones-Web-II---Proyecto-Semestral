package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarPublicacion atiende GET /api/v1/publicaciones.
func (s *Server) ListarPublicacion(w http.ResponseWriter, _ *http.Request) {
	publicaciones := s.Storage.ListarPublicacion()
	RespondJSON(w, http.StatusOK, publicaciones)
}

// ObtenerPublicacion atiende GET /api/v1/publicaciones/{id}.
func (s *Server) ObtenerPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	publicacion, encontrado := s.Storage.BuscarPublicacionPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Publicación no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, publicacion)
}

// CrearPublicacion atiende POST /api/v1/publicaciones.
func (s *Server) CrearPublicacion(w http.ResponseWriter, r *http.Request) {
	var nueva models.Publicacion
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if strings.TrimSpace(nueva.Titulo) == "" {
		RespondError(w, http.StatusBadRequest, "El campo título es obligatorio")
		return
	}

	creada := s.Storage.CrearPublicacion(nueva)
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarPublicacion atiende PUT /api/v1/publicaciones/{id}.
func (s *Server) ActualizarPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Publicacion
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	if strings.TrimSpace(datos.Titulo) == "" {
		RespondError(w, http.StatusBadRequest, "El campo título es obligatorio")
		return
	}

	actualizada, encontrada := s.Storage.ActualizarPublicacion(id, datos)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "Publicación no encontrada")
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarPublicacion atiende DELETE /api/v1/publicaciones/{id}.
func (s *Server) BorrarPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if !s.Storage.BorrarPublicacion(id) {
		RespondError(w, http.StatusNotFound, "Publicación no encontrada")
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
