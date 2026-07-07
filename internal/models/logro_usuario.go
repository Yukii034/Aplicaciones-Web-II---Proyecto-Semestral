package models

import "time"

type Logro_Usuario struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	FechaDesbl time.Time `json:"fechas_desbloqueado" gorm:"not null"`
	UsuarioID  int       `json:"usuario_id"`
	LogroID    int       `json:"logro_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Logro      Logro     `json:"logro,omitempty" gorm:"foreignKey:LogroID"`
	Usuario    Usuario   `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}
