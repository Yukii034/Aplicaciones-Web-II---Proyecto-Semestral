package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

func (s *Server) ListarUsuarios(w http.ResponseWriter, _ *http.Request) {
	usuarios := s.Usuario.ListarUsuarios()
	RespondJSON(w, http.StatusOK, usuarios)
}

func (s *Server) ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	usuarios := s.Usuario.ListarUsuarios()
	var usuario models.Usuario
	encontrado := false
	for _, u := range usuarios {
		if u.ID == id {
			usuario = u
			encontrado = true
			break
		}
	}
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, usuario)
}

func (s *Server) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		RespondError(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
		return
	}
	s.Usuario.CrearUsuario(usuario)
	RespondJSON(w, http.StatusCreated, usuario)
}

func (s *Server) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		RespondError(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
		return
	}
	usuario.ID = id
	s.Usuario.ActualizarUsuario(id, usuario)
	RespondJSON(w, http.StatusOK, usuario)
}

func (s *Server) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	if err := s.Usuario.BorrarUsuario(id); err != nil {
		RespondError(w, http.StatusNotFound, "Usuario no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"message": "Usuario eliminado correctamente"})
}
