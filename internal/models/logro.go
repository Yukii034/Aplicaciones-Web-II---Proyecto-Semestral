package models

import "time"

type Logro struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	Nombre           string    `json:"nombre" gorm:"not null"`
	Descripcion      string    `json:"descripcion"`
	PuntosRequeridos float64   `json:"puntos_requeridos"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
