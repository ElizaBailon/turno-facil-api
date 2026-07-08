package repository

import (
	"turno-facil-api/internal/models"

	"gorm.io/gorm"
)

// ServiciosRepository define el contrato para el catálogo de trabajos del taller
type ServiciosRepository interface {
	Crear(s *models.Servicio) error
	Listar() ([]models.Servicio, error)
	ObtenerPorID(id int) (models.Servicio, error)
}

type sqliteServiciosRepo struct {
	db *gorm.DB
}

// NewServiciosRepository ahora retorna la interfaz
func NewServiciosRepository(db *gorm.DB) ServiciosRepository {
	return &sqliteServiciosRepo{db: db}
}

func (r *sqliteServiciosRepo) Crear(s *models.Servicio) error {
	return r.db.Create(s).Error
}

func (r *sqliteServiciosRepo) Listar() ([]models.Servicio, error) {
	var servicios []models.Servicio
	err := r.db.Find(&servicios).Error
	return servicios, err
}

func (r *sqliteServiciosRepo) ObtenerPorID(id int) (models.Servicio, error) {
	var servicio models.Servicio
	err := r.db.First(&servicio, id).Error
	return servicio, err
}
