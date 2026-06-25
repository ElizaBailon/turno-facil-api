package repository

import (
	"turno-facil-api/internal/models"

	"gorm.io/gorm"
)

// MecanicosRepository define las operaciones permitidas para el personal técnico
type MecanicosRepository interface {
	Crear(m *models.Mecanico) error
	Listar() ([]models.Mecanico, error)
	ObtenerPorID(id int) (models.Mecanico, error)
	// 🌟 Agregamos la firma del método que solicita el servicio para cumplir las reglas del Hito 3
	CambiarDisponibilidad(id int, activo bool) error
}

type sqliteMecanicosRepo struct {
	db *gorm.DB
}

// NewMecanicosRepository ahora retorna la interfaz desacoplada
func NewMecanicosRepository(db *gorm.DB) MecanicosRepository {
	return &sqliteMecanicosRepo{db: db}
}

func (r *sqliteMecanicosRepo) Crear(m *models.Mecanico) error {
	return r.db.Create(m).Error
}

func (r *sqliteMecanicosRepo) Listar() ([]models.Mecanico, error) {
	var mecanicos []models.Mecanico
	err := r.db.Find(&mecanicos).Error
	return mecanicos, err
}

func (r *sqliteMecanicosRepo) ObtenerPorID(id int) (models.Mecanico, error) {
	var mecanico models.Mecanico
	err := r.db.First(&mecanico, id).Error
	return mecanico, err
}

// 🌟 Implementación del método utilizando GORM para actualizar SQLite
func (r *sqliteMecanicosRepo) CambiarDisponibilidad(id int, activo bool) error {
	// Busca en la tabla correspondinte al modelo Mecanico y actualiza la columna 'activo'
	return r.db.Model(&models.Mecanico{}).Where("id = ?", id).Update("activo", activo).Error
}
