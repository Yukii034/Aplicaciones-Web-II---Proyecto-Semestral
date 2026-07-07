package models

import "time"

type Calificacion struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Comentarios string    `json:"comentarios" gorm:"not null"`
	UsuarioID   int       `json:"usuario_id"`
	AcuerdoID   int       `json:"acuerdo_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Acuerdo     Acuerdo   `json:"acuerdo,omitempty" gorm:"foreignKey:AcuerdoID"`
	Usuario     Usuario   `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}
