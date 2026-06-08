package models

import "time"

type Inventario struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	Nombre         string    `json:"nombre" gorm:"not null"`
	Descripcion    string    `json:"descripcion"`
	Categoria      string    `json:"categoria"`
	EstadoObjeto   string    `json:"estado_objeto"`
	Disponibilidad string    `json:"disponibilidad"`
	Cantidad       int       `json:"cantidad"`
	FotoUrl        string    `json:"foto_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UsuarioID      int       `json:"usuario_id"`
}
