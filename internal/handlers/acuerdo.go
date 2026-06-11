package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

func (s *Server) ListarAcuerdo(w http.ResponseWriter, _ *http.Request) {
	acuerdos := s.Storage.ListarAcuerdos()
	RespondJSON(w, http.StatusOK, acuerdos)
}

func (s *Server) ObtenerAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	acuerdo, encontrado := s.Storage.BuscarAcuerdoPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "Acuerdo no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, acuerdo)
}

func (s *Server) CrearAcuerdo(w http.ResponseWriter, r *http.Request) {
	var acuerdo models.Acuerdo
	if err := json.NewDecoder(r.Body).Decode(&acuerdo); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	s.Storage.CrearAcuerdo(acuerdo)
	RespondJSON(w, http.StatusCreated, acuerdo)
}

func (s *Server) ActualizarAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	var acuerdo models.Acuerdo
	if err := json.NewDecoder(r.Body).Decode(&acuerdo); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	updatedAcuerdo, ok := s.Storage.ActualizarAcuerdo(id, acuerdo)
	if !ok {
		RespondError(w, http.StatusNotFound, "Acuerdo no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, updatedAcuerdo)
}

func (s *Server) EliminarAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	if !s.Storage.BorrarAcuerdo(id) {
		RespondError(w, http.StatusNotFound, "Acuerdo no encontrado")
		return
	}
}
