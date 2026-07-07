package models

import "time"

type Reputacion struct {
	ID                   int       `json:"id" gorm:"primaryKey"`
	PuntosTotales        int       `json:"puntos_totales" gorm:"not null"`
	Nivel                int       `json:"nivel"`
	AcuerdosCompl        int       `json:"acuerdo_compl"`
	CalificacionPromedio float64   `json:"calificacion_pro"`
	UsuarioID            int       `json:"usuario_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	Usuario              Usuario   `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}
