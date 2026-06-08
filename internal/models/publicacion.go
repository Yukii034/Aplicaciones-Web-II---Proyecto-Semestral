package models

import "time"

type Publicacion struct {
	ID                int       `json:"id" gorm:"primaryKey"`
	Titulo            string    `json:"titulo" gorm:"not null"`
	TipoOferta        string    `json:"tipo_oferta"`
	EstadoPublicacion string    `json:"estado_publicacion"`
	Mensaje           string    `json:"mensaje"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UsuarioID         int       `json:"usuario_id"`
	InventarioID      int       `json:"inventario_id"`
}
