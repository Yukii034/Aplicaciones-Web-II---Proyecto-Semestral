package models

import "time"

type Publicacion struct {
	ID                 int       `json:"id" gorm:"primaryKey"`
	Titulo             string    `json:"titulo" gorm:"not null"`
	Tipo_Oferta        string    `json:"tipo_oferta"`
	Estado_Publicacion int       `json:"stock"`
	Mensaje            string    `json:"mensaje"`
	Created_at         time.Time `json:"created_at"`
	Updated_at         time.Time `json:"updated_at"`
	UsuarioID          int       `json:"usuario_id"`
	InventarioID       int       `json:"inventario_id"`
}
