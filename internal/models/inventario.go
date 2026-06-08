package models

import "time"

type Inventario struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	Nombre         string    `json:"nombre" gorm:"not null"`
	Descripcion    string    `json:"descripcion"`
	Categoria      string    `json:"categoria"`
	Estado_Objeto  string    `json:"estado_objeto"`
	Disponibilidad string    `json:"disponibilidad"`
	Cantidad       int       `json:"cantidad"`
	Foto_Url       string    `json:"foto_url"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
	UsuarioID      int       `json:"usuario_id"`
}
