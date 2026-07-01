package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto-semestral/internal/models"
)

func (s *Server) ListarAcuerdoItem(w http.ResponseWriter, _ *http.Request) {
	items := s.AcuerdoItem.ListarAcuerdoItems()
	RespondJSON(w, http.StatusOK, items)
}

func (s *Server) ObtenerAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r) // cambiado por el params.go
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	items := s.AcuerdoItem.ListarAcuerdoItems()
	for _, item := range items {
		if item.ID == id {
			RespondJSON(w, http.StatusOK, item)
			return
		}
	}
	RespondError(w, http.StatusNotFound, "AcuerdoItem no encontrado")
}

func (s *Server) CrearAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	var item models.AcuerdoItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	s.AcuerdoItem.CrearAcuerdoItem(item)
	RespondJSON(w, http.StatusCreated, item)
}

func (s *Server) ActualizarAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	var item models.AcuerdoItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	updatedItem, err := s.AcuerdoItem.ActualizarAcuerdoItem(id, item)
	if err != nil {
		RespondError(w, http.StatusNotFound, "AcuerdoItem no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, updatedItem)
}

func (s *Server) EliminarAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	if err := s.AcuerdoItem.BorrarAcuerdoItem(id); err != nil {
		RespondError(w, http.StatusNotFound, "AcuerdoItem no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"message": "AcuerdoItem eliminado correctamente"})
}
