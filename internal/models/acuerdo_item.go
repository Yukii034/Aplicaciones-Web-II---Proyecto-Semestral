package models

//Struct para el item del acuerdo

type AcuerdoItem struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	AcuerdoID    int    `json:"acuerdo_id"`
	InventarioID int    `json:"inventario_id"`
	Rol          string `json:"rol"`
}
