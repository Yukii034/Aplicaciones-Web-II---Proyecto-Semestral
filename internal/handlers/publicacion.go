package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"proyecto-semestral/internal/models"
)

// ListarPublicacion atiende GET /api/v1/publicaciones.
func (s *Server) ListarPublicacion(w http.ResponseWriter, _ *http.Request) {
	publicaciones := s.Publicacion.ListarPublicacion() // cambiar a que todos tengan su servicio con su handler
	RespondJSON(w, http.StatusOK, publicaciones)
}

// ObtenerPublicacion atiende GET /api/v1/publicaciones/{id}.
func (s *Server) ObtenerPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	publicacion, err := s.Publicacion.BuscarPublicacion(id)
	if err != nil { // agg los errores ya que ahora retornan errores
		RespondError(w, statusDeError(err), err.Error())
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

	creada, err := s.Publicacion.CrearPublicacion(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarPublicacion atiende PUT /api/v1/publicaciones/{id}.
func (s *Server) ActualizarPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
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

	actualizada, err := s.Publicacion.ActualizarPublicacion(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarPublicacion atiende DELETE /api/v1/publicaciones/{id}.
func (s *Server) BorrarPublicacion(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Publicacion.BorrarPublicacion(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
