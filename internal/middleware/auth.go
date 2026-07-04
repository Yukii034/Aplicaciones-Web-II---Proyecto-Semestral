package middleware

import (
	"net/http"

	"proyecto-semestral/internal/service"

	"context"
	"strings"
)

type claveContext string

const ClaveUsuarioID claveContext = "usuarioID"
const ClaveTipo claveContext = "tipo" // para roles

// Esto se denomina factory de middleware
func Auth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encabezado := r.Header.Get("Authorization")
			partes := strings.SplitN(encabezado, " ", 2)
			if len(partes) != 2 || partes[0] != "Bearer" {
				responderNoAutorizado(w)
				return
			}
			usuarioID, tipo, err := auth.ValidarToken(partes[1])
			if err != nil {
				responderNoAutorizado(w)
				return
			}
			ctx := context.WithValue(r.Context(), ClaveUsuarioID, usuarioID)
			ctx = context.WithValue(ctx, ClaveTipo, tipo) // para roles
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func responderNoAutorizado(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error": "Token inexistente o invalido"}`))
}

// para el rol de admin
func SoloAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tipo, ok := r.Context().Value(ClaveTipo).(string)
		if !ok || tipo != "admin" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "acceso denegado, se requiere rol admin"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}
