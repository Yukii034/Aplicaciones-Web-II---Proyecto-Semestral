package models

//Struct de usuario, con sus respectivos campos y etiquetas para JSON y GORM

type Usuario struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	Nombre          string `json:"nombre" gorm:"not null"`
	Email           string `json:"email" gorm:"unique;not null"`
	Contraseña      string `json:"contraseña" gorm:"not null"`
	Tipo_de_Usuario string `json:"tipo_usuario" gorm:"not null"`
	Ciudad          string `json:"ciudad"`
	Telefono        string `json:"telefono"`
	Reputacion      string `json:"reputacion" gorm:"foreignKey:UsuarioID"`
}
