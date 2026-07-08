package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto-semestral/internal/models"
)

// ListarAcuerdo atiende GET /api/v1/acuerdos.
func (s *Server) ListarAcuerdo(w http.ResponseWriter, _ *http.Request) {
	acuerdos := s.Acuerdo.ListarAcuerdos()
	RespondJSON(w, http.StatusOK, acuerdos)
}

// ObtenerAcuerdo atiende GET /api/v1/acuerdos/{id}.
func (s *Server) ObtenerAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	acuerdo, err := s.Acuerdo.BuscarAcuerdo(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, acuerdo)
}

// CrearAcuerdo atiende POST /api/v1/acuerdos.
func (s *Server) CrearAcuerdo(w http.ResponseWriter, r *http.Request) {
	var nuevo models.Acuerdo
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	creado, err := s.Acuerdo.CrearAcuerdo(nuevo)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarAcuerdo atiende PUT /api/v1/acuerdos/{id}.
func (s *Server) ActualizarAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	var acuerdo models.Acuerdo
	if err := json.NewDecoder(r.Body).Decode(&acuerdo); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	actualizado, err := s.Acuerdo.ActualizarAcuerdo(id, acuerdo)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizado)
}

// EliminarAcuerdo atiende DELETE /api/v1/acuerdos/{id}.
func (s *Server) EliminarAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.Acuerdo.BorrarAcuerdo(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "Acuerdo eliminado correctamente"})
}
