package storage

import (
	"time"

	"proyecto-semestral/internal/models"
)

// SembrarSiVacio inserta datos de ejemplo solo si no hay usuarios todavia.
// Evita duplicar datos cada vez que se llama (ej. al reiniciar el contenedor).
func (a *AlmacenSQLite) SembrarSiVacio() {
	if len(a.ListarUsuarios()) > 0 {
		return // ya hay datos, no hacer nada
	}

	// --- Usuarios ---
	u1 := a.CrearUsuario(models.Usuario{
		Nombre:       "Admin",
		Email:        "admin@test.com",
		PasswordHash: "hash_ejemplo",
		Tipo:         "admin",
		Ciudad:       "Manta",
		Telefono:     "+593999999999",
		Reputacion:   "buena",
	})

	u2 := a.CrearUsuario(models.Usuario{
		Nombre:       "Maria Perez",
		Email:        "maria@test.com",
		PasswordHash: "hash_ejemplo2",
		Tipo:         "cliente",
		Ciudad:       "Portoviejo",
		Telefono:     "+593988888888",
		Reputacion:   "buena",
	})

	// --- Inventario ---
	inv := a.CrearInventario(models.Inventario{
		Nombre:         "Silla de escritorio",
		Descripcion:    "Silla en buen estado",
		Categoria:      "Muebles",
		EstadoObjeto:   "usado",
		Disponibilidad: "disponible",
		Cantidad:       1,
		UsuarioID:      u1.ID,
	})

	// --- Publicacion ---
	pub := a.CrearPublicacion(models.Publicacion{
		Titulo:            "Cambio silla de escritorio",
		TipoOferta:        "intercambio",
		EstadoPublicacion: "activa",
		Mensaje:           "Cambio por algo util para estudiar",
		UsuarioID:         u1.ID,
		InventarioID:      inv.ID,
	})

	// --- Reputacion ---
	a.CrearReputacion(models.Reputacion{
		PuntosTotales:        50,
		Nivel:                1,
		AcuerdosCompl:        0,
		CalificacionPromedio: 0,
		UsuarioID:            u1.ID,
	})

	a.CrearReputacion(models.Reputacion{
		PuntosTotales:        10,
		Nivel:                1,
		AcuerdosCompl:        0,
		CalificacionPromedio: 0,
		UsuarioID:            u2.ID,
	})

	// --- Logro ---
	logro := a.CrearLogro(models.Logro{
		Nombre:           "Primer intercambio",
		Descripcion:      "Completa tu primer acuerdo",
		PuntosRequeridos: 10,
	})

	// --- Logro_Usuario ---
	a.CrearLogro_Usuario(models.Logro_Usuario{
		FechaDesbl: time.Now(),
		UsuarioID:  u1.ID,
		LogroID:    logro.ID,
	})

	// --- Acuerdo ---
	acuerdo := a.CrearAcuerdo(models.Acuerdo{
		PublicacionID:            pub.ID,
		IDOfertante:              u2.ID,
		IDPublicador:             u1.ID,
		Tipo:                     "intercambio",
		Estado:                   "pendiente",
		Mensaje_Inicial:          "Me interesa la silla",
		Motivo_Cancelacion:       "",
		Confirmacion_Solicitante: false,
	})

	// --- Acuerdo_Item ---
	a.CrearAcuerdoItem(models.AcuerdoItem{
		AcuerdoID:    acuerdo.ID,
		InventarioID: inv.ID,
		Rol:          "ofrecido",
	})

	// --- Calificacion ---
	a.CrearCalificacion(models.Calificacion{
		Comentarios: "Buena experiencia, todo tal como se acordo",
		UsuarioID:   u2.ID,
		AcuerdoID:   acuerdo.ID,
	})
}
