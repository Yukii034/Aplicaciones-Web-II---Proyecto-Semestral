package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto-semestral/internal/models"
)

func (s *Server) ListarAcuerdo(w http.ResponseWriter, _ *http.Request) {
	acuerdos := s.Acuerdo.ListarAcuerdos()
	RespondJSON(w, http.StatusOK, acuerdos)
}

func (s *Server) ObtenerAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	// Fallback: buscar en la lista de acuerdos si no existe un método ObtenerAcuerdo
	acuerdos := s.Acuerdo.ListarAcuerdos()
	var acuerdo models.Acuerdo
	encontrado := false
	for _, a := range acuerdos {
		if a.ID == id {
			acuerdo = a
			encontrado = true
			break
		}
	}
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
	s.Acuerdo.CrearAcuerdo(acuerdo)
	RespondJSON(w, http.StatusCreated, acuerdo)
}

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
	updatedAcuerdo, err := s.Acuerdo.ActualizarAcuerdo(id, acuerdo)
	if err != nil {
		RespondError(w, http.StatusNotFound, "Acuerdo no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, updatedAcuerdo)
}

func (s *Server) EliminarAcuerdo(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	if err := s.Acuerdo.BorrarAcuerdo(id); err != nil {
		RespondError(w, http.StatusNotFound, "Acuerdo no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"message": "Acuerdo eliminado correctamente"})
}
