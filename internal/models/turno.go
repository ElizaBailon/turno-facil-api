package models

import "time"

type Turno struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	VehiculoID  int       `json:"vehiculo_id" gorm:"not null"`
	Vehiculo    Vehiculo  `json:"vehiculo,omitempty" gorm:"foreignKey:VehiculoID"`
	MecanicoID  int       `json:"mecanico_id" gorm:"not null"`
	Mecanico    Mecanico  `json:"mecanico,omitempty" gorm:"foreignKey:MecanicoID"`
	ServicioID  int       `json:"servicio_id" gorm:"not null"`
	Servicio    Servicio  `json:"servicio,omitempty" gorm:"foreignKey:ServicioID"`
	FechaHora   time.Time `json:"fecha_hora" gorm:"not null"`
	DuracionEst int       `json:"duracion_estimada"`
	Estado      string    `json:"estado" gorm:"default:'pendiente'"`
	Notas       string    `json:"notas"`
}
