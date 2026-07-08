package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto-semestral/internal/models"
)

// ListarAcuerdoItem atiende GET /api/v1/acuerdo-items.
func (s *Server) ListarAcuerdoItem(w http.ResponseWriter, _ *http.Request) {
	items := s.AcuerdoItem.ListarAcuerdoItems()
	RespondJSON(w, http.StatusOK, items)
}

// ObtenerAcuerdoItem atiende GET /api/v1/acuerdo-items/{id}.
func (s *Server) ObtenerAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r) // cambiado por el params.go
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	item, err := s.AcuerdoItem.BuscarAcuerdoItem(id)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, item)
}

// CrearAcuerdoItem atiende POST /api/v1/acuerdo-items.
func (s *Server) CrearAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	var nuevo models.AcuerdoItem
	if err := json.NewDecoder(r.Body).Decode(&nuevo); err != nil {
		RespondError(w, http.StatusBadRequest, "Datos inválidos: "+err.Error())
		return
	}

	creado, err := s.AcuerdoItem.CrearAcuerdoItem(nuevo)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, creado)
}

// ActualizarAcuerdoItem atiende PUT /api/v1/acuerdo-items/{id}.
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

	actualizado, err := s.AcuerdoItem.ActualizarAcuerdoItem(id, item)
	if err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, actualizado)
}

// EliminarAcuerdoItem atiende DELETE /api/v1/acuerdo-items/{id}.
func (s *Server) EliminarAcuerdoItem(w http.ResponseWriter, r *http.Request) {
	id, err := idDeURL(r)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "El ID debe ser un número entero")
		return
	}

	if err := s.AcuerdoItem.BorrarAcuerdoItem(id); err != nil {
		RespondError(w, statusDeError(err), err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, map[string]string{"message": "AcuerdoItem eliminado correctamente"})
}
