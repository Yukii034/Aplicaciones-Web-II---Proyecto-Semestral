package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto-semestral/internal/models"
)

// ListarUsuarios atiende GET /api/v1/usuarios.
func (s *Server) ListarUsuarios(w http.ResponseWriter, _ *http.Request) {
	usuarios := s.Usuario.ListarUsuarios()
	RespondJSON(w, http.StatusOK, usuarios)
}

// ObtenerUsuario atiende GET /api/v1/usuarios/{id}.
func (s *Server) ObtenerUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	usuario, err := s.Usuario.BuscarUsuario(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, usuario)
}

// CrearUsuario atiende POST /api/v1/usuarios.
func (s *Server) CrearUsuario(w http.ResponseWriter, r *http.Request) {
	var nuevo models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		RespondError(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
		return
	}

	creado, err := s.Usuario.CrearUsuario(nuevo)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarUsuario atiende PUT /api/v1/usuarios/{id}.
func (s *Server) ActualizarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var datos models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		RespondError(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la solicitud")
		return
	}

	actualizado, err := s.Usuario.ActualizarUsuario(id, datos)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizado)
}

// EliminarUsuario atiende DELETE /api/v1/usuarios/{id}.
func (s *Server) EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Usuario.BorrarUsuario(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "Usuario eliminado correctamente"})
}
