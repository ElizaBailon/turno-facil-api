package models

type Cliente struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Cedula   string `json:"cedula" gorm:"not null;unique"`
	Nombre   string `json:"nombre" gorm:"not null"`
	Telefono string `json:"telefono"`
	Email    string `json:"email"`
	// 🌟 Agregamos constraint:OnDelete:CASCADE para que borre sus vehículos automáticamente
	Vehiculo []Vehiculo `json:"vehiculo" gorm:"foreignKey:ClienteID;constraint:OnDelete:CASCADE;"`
}
