package repository

import (
	"turno-facil-api/internal/models"

	"gorm.io/gorm"
)

type UsuariosRepository interface {
	Crear(u *models.Usuario) error
	ObtenerPorUsername(username string) (models.Usuario, error)
}

type sqliteUsuariosRepo struct {
	db *gorm.DB
}

func NewUsuariosRepository(db *gorm.DB) UsuariosRepository {
	return &sqliteUsuariosRepo{db: db}
}

func (r *sqliteUsuariosRepo) Crear(u *models.Usuario) error {
	return r.db.Create(u).Error
}

func (r *sqliteUsuariosRepo) ObtenerPorUsername(username string) (models.Usuario, error) {
	var usuario models.Usuario
	err := r.db.Where("username = ?", username).First(&usuario).Error
	return usuario, err
}
