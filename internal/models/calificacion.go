package models

import "time"

type Calificacion struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Comentarios string    `json:"puntos_totales" gorm:"not null"`
	UsuarioID   int       `json:"usuario_id"`
	AcuerdoID   int       `json:"acuerdo_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
