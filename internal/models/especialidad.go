package models

type Especialidad struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Nombre      string `json:"nombre" gorm:"not null;unique"` // Ej: "Frenos", "Motores"
	Descripcion string `json:"descripcion"`
}
