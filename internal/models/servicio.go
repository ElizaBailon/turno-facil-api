package models

type Servicio struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	Nombre       string  `json:"nombre" gorm:"not null;unique"`
	DuracionMins int     `json:"duracion_mins" gorm:"not null"` // Importante para el algoritmo
	PrecioEst    float64 `json:"precio_estimado"`
}
