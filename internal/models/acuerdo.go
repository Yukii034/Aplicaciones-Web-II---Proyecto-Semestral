package models

import "time"

//Struct de Acuerdo con los campos necesarios para la creación de un acuerdo

type Acuerdo struct {
	ID                       int         `json:"id" gorm:"primaryKey"`
	PublicacionID            int         `json:"publicacion_id" gorm:"not null"`
	IDOfertante              int         `json:"id_ofertante" gorm:"not null"`
	IDSolicitante            int         `json:"id_solicitante" gorm:"not null"`
	Tipo                     string      `json:"tipo" gorm:"not null"`
	Estado                   string      `json:"estado"`
	Mensaje_Inicial          string      `json:"mensaje_inicial"`
	Motivo_Cancelacion       string      `json:"motivo_cancelacion"`
	Confirmacion_Solicitante bool        `json:"confirmacion_solicitante"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
	Publicacion              Publicacion `json:"publicacion,omitempty" gorm:"foreignKey:PublicacionID"`
	Ofertante                Usuario     `json:"ofertante,omitempty" gorm:"foreignKey:IDOfertante"`
	Solicitante              Usuario     `json:"solicitante,omitempty" gorm:"foreignKey:IDSolicitante"`
}
