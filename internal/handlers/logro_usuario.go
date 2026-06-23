package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

// ListarLogroUsuario atiende GET /api/v1/logro_usuarios.
func (s *Server) ListarLogro_Usuario(w http.ResponseWriter, _ *http.Request) {
	logro_usuarios := s.Logro_Usuario.ListarLogro_Usuario()
	RespondJSON(w, http.StatusOK, logro_usuarios)
}

// ObtenerLogro_Usuario atiende GET /api/v1/logro_usuarios/{id}.
func (s *Server) ObtenerLogro_Usuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	logro_usuario, err := s.Logro_Usuario.BuscarLogro_Usuario(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, logro_usuario)
}

// CrearLogro_Usuario atiende POST /api/v1/logro_usuarios.
func (s *Server) CrearLogro_Usuario(w http.ResponseWriter, r *http.Request) {
	var nueva models.Logro_Usuario
	if err := json.NewDecoder(r.Body).Decode(&nueva); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	creada, err := s.Logro_Usuario.CrearLogro_Usuario(nueva)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
	}
	RespondJSON(w, http.StatusCreated, creada)
}

// ActualizarLogro_Usuario atiende PUT /api/v1/logro_usuarios/{id}.
func (s *Server) ActualizarLogro_Usuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Logro_Usuario
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}

	actualizada, err := s.Logro_Usuario.ActualizarLogro_Usuario(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizada)
}

// BorrarLogro_Usuario atiende DELETE /api/v1/logro_usuarios/{id}.
func (s *Server) BorrarLogro_Usuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Logro_Usuario.BorrarLogro_Usuario(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
