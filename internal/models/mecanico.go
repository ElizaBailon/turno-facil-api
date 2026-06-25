package models

import "time"

type Mecanico struct {
	ID             int          `json:"id" gorm:"primaryKey"`
	Nombre         string       `json:"nombre" gorm:"not null"`
	EspecialidadID int          `json:"especialidad_id" gorm:"not null"` // Llave foránea
	Especialidad   Especialidad `json:"especialidad,omitempty" gorm:"foreignKey:EspecialidadID"`
	Activo         bool         `json:"activo" gorm:"default:true"`
	CreatedAt      time.Time    `json:"created_at"`
}
