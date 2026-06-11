package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"proyecto-semestral/internal/models"
)

func (s *Server) ListarAcuerdoItem(w http.ResponseWriter, _ *http.Request) {
	items := s.Storage.ListarAcuerdoItems()
	RespondJSON(w, http.StatusOK, items)
}

func (s *Server) ObtenerAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	item, encontrado := s.Storage.BuscarAcuerdoItemPorID(id)
	if !encontrado {
		RespondError(w, http.StatusNotFound, "AcuerdoItem no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, item)
}

func (s *Server) CrearAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	var item models.AcuerdoItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	s.Storage.CrearAcuerdoItem(item)
	RespondJSON(w, http.StatusCreated, item)
}

func (s *Server) ActualizarAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	var item models.AcuerdoItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}
	updatedItem, ok := s.Storage.ActualizarAcuerdoItem(id, item)
	if !ok {
		RespondError(w, http.StatusNotFound, "AcuerdoItem no encontrado")
		return
	}
	RespondJSON(w, http.StatusOK, updatedItem)
}

func (s *Server) EliminarAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}
	if !s.Storage.BorrarAcuerdoItem(id) {
		RespondError(w, http.StatusNotFound, "AcuerdoItem no encontrado")
		return
	}
}
