package models

//Struct para el item del acuerdo

type AcuerdoItem struct {
	ID           int        `json:"id" gorm:"primaryKey"`
	AcuerdoID    int        `json:"acuerdo_id" gorm:"not null"`
	InventarioID int        `json:"inventario_id" gorm:"not null"`
	Rol          string     `json:"rol" gorm:"not null"`
	Acuerdo      Acuerdo    `json:"acuerdo,omitempty" gorm:"foreignKey:AcuerdoID"`
	Inventario   Inventario `json:"inventario,omitempty" gorm:"foreignKey:InventarioID"`
}
