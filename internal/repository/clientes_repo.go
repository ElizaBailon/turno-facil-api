package repository

import (
	"errors"
	"strings"
	"turno-facil-api/internal/models"

	"gorm.io/gorm"
)

// ClientesRepository define las operaciones que cualquier base de datos de clientes debe cumplir
type ClientesRepository interface {
	Crear(c *models.Cliente) error
	BuscarPorCedula(cedula string) (models.Cliente, error)
	Actualizar(c *models.Cliente) error
	Eliminar(id int) error // 🌟 CORREGIDO: Cambiado de string a int para mantener la consistencia del hito
}

// sqliteClientesRepo es la implementación real usando GORM y SQLite
type sqliteClientesRepo struct {
	db *gorm.DB
}

// NewClientesRepository ahora devuelve la interfaz, no la estructura concreta
func NewClientesRepository(db *gorm.DB) ClientesRepository {
	return &sqliteClientesRepo{db: db}
}

func (r *sqliteClientesRepo) Crear(c *models.Cliente) error {
	err := r.db.Create(c).Error
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
			return errors.New("duplicado")
		}
		return err
	}
	return nil
}

func (r *sqliteClientesRepo) BuscarPorCedula(cedula string) (models.Cliente, error) {
	var cliente models.Cliente
	// 🌟 Trae los carros vinculados automáticamente gracias al Preload en plural
	err := r.db.Preload("Vehiculo").Where("cedula = ?", cedula).First(&cliente).Error
	return cliente, err
}

func (r *sqliteClientesRepo) Actualizar(c *models.Cliente) error {
	// 🌟 FullSaveAssociations guarda el cliente y procesa sus vehículos vinculados en cascada
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(c).Error
}

// 🌟 NUEVO: Implementación física del método Eliminar exigido por la interfaz
func (r *sqliteClientesRepo) Eliminar(id int) error {
	res := r.db.Delete(&models.Cliente{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("cliente no encontrado")
	}
	return nil
}
