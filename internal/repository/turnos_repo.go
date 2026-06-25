package repository

import (
	"turno-facil-api/internal/models"

	"gorm.io/gorm"
)

// TurnosRepository define el contrato para la persistencia de turnos
type TurnosRepository interface {
	Crear(t *models.Turno) error
	Listar() ([]models.Turno, error)
	ListarPorMecanico(mecanicoID int) ([]models.Turno, error)
	Eliminar(id int) error            // 🌟 ¡Añade esta línea!
	Actualizar(t *models.Turno) error // 🌟 ¡Añade esta línea!
	BuscarPorID(id int) (models.Turno, error)
}

type sqliteTurnosRepo struct {
	db *gorm.DB
}

func NewTurnosRepository(db *gorm.DB) TurnosRepository {
	return &sqliteTurnosRepo{db: db}
}

func (r *sqliteTurnosRepo) Crear(t *models.Turno) error {
	// 1. Primero insertamos el turno de forma normal en SQLite
	if err := r.db.Create(t).Error; err != nil {
		return err
	}

	// 2. 🌟 Cargamos las relaciones para rellenar el struct antes de que regrese al controlador
	return r.db.Preload("Vehiculo.Cliente").
		Preload("Mecanico.Especialidad").
		Preload("Servicio").
		First(t, t.ID).Error
}

func (r *sqliteTurnosRepo) Listar() ([]models.Turno, error) {
	var turnos []models.Turno
	err := r.db.Preload("Servicio").
		Preload("Vehiculo.Cliente").
		Preload("Mecanico.Especialidad"). // 🌟 Esto llena la especialidad
		Find(&turnos).Error
	return turnos, err
}

func (r *sqliteTurnosRepo) ListarPorMecanico(mecanicoID int) ([]models.Turno, error) {
	var turnos []models.Turno
	err := r.db.Preload("Servicio").
		Preload("Vehiculo.Cliente").
		Preload("Mecanico.Especialidad").
		Where("mecanico_id = ?", mecanicoID).
		Find(&turnos).Error
	return turnos, err
}

func (r *sqliteTurnosRepo) Eliminar(id int) error {
	return r.db.Delete(&models.Turno{}, id).Error
}

func (r *sqliteTurnosRepo) Actualizar(t *models.Turno) error {
	return r.db.Save(t).Error
}

func (r *sqliteTurnosRepo) BuscarPorID(id int) (models.Turno, error) {
	var turno models.Turno
	// 🌟 Aquí ocurre la magia de traer los objetos relacionados
	err := r.db.Preload("Servicio").
		Preload("Vehiculo.Cliente").
		Preload("Mecanico").
		Preload("Mecanico.Especialidad").
		First(&turno, id).Error

	if err != nil {
		return models.Turno{}, err
	}
	return turno, nil
}
