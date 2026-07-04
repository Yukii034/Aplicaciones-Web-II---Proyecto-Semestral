package handlers

import (
	"encoding/json"
	"net/http"
)

type credenciales struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Tipo     string `json:"tipo"` // para roles
}

func (s *Server) Registrar(w http.ResponseWriter, r *http.Request) {
	var creds credenciales
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	usuario, err := s.Auth.Registrar(creds.Email, creds.Password, creds.Tipo)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondJSON(w, http.StatusCreated, usuario)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var creds credenciales
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		RespondError(w, http.StatusBadRequest, "JSON inválido: "+err.Error())
		return
	}
	token, err := s.Auth.Login(creds.Email, creds.Password)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondJSON(w, http.StatusOK, map[string]string{"token": token})
}
