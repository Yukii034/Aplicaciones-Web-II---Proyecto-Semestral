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
	logro_usuarios := s.Storage.ListarLogro_Usuario()
	RespondJSON(w, http.StatusOK, logro_usuarios)
}

// ObtenerLogro_Usuario atiende GET /api/v1/logro_usuarios/{id}.
func (s *Server) ObtenerLogro_Usuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	logro_usuario, encontrado := s.Storage.BuscarLogro_UsuarioPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "logro usuario no encontrada")
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
	if nueva.UsuarioID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	if nueva.LogroID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	creada := s.Storage.CrearLogro_Usuario(nueva)
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

	if datos.UsuarioID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo usuario_id es obligatorio")
		return
	}

	if datos.LogroID == 0 {
		RespondError(w, http.StatusBadRequest, "El campo logro_id es obligatorio")
		return
	}

	actualizada, encontrada := s.Storage.ActualizarLogro_Usuario(id, datos)
	if !encontrada {
		RespondError(w, http.StatusNotFound, "logro usuario no encontrado")
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

	if !s.Storage.BorrarLogro_Usuario(id) {
		RespondError(w, http.StatusNotFound, "logro usuario no encontrada")
		return
	}

	RespondJSON(w, http.StatusNoContent, nil)
}
