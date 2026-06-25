package models

type Vehiculo struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Placa     string `json:"placa" gorm:"not null;unique"` // Ej: "MAN-1234"
	Marca     string `json:"marca" gorm:"not null"`        // Ej: "Chevrolet"
	Modelo    string `json:"modelo"`                       // Ej: "Sail"
	ClienteID int    `json:"cliente_id" gorm:"not null"`   // Llave foránea
	// 🌟 Cambia la etiqueta JSON para que ignore este campo si viene vacío o al serializar
	Cliente *Cliente `json:"cliente,omitempty" gorm:"foreignKey:ClienteID"`
}
