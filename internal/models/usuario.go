package models

import "time"

type Usuario struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Username  string     `gorm:"uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"not null" json:"-"`             // No se envía en las respuestas HTTP
	Rol       string     `gorm:"default:'mecanico'" json:"rol"` // Ej: "admin", "mecanico"
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// Estructuras auxiliares para las peticiones HTTP
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Rol   string `json:"rol"`
}
