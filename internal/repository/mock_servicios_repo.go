package repository

import (
	"turno-facil-api/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockServiciosRepository struct {
	mock.Mock
}

func (m *MockServiciosRepository) Listar() ([]models.Servicio, error) {
	args := m.Called()
	return args.Get(0).([]models.Servicio), args.Error(1)
}

func (m *MockServiciosRepository) Crear(servicio *models.Servicio) error {
	args := m.Called(servicio)
	return args.Error(0)
}

func (m *MockServiciosRepository) ObtenerPorID(id int) (models.Servicio, error) { // O la estructura que use tu interfaz
	// Si tu método real devuelve models.Servicio:
	args := m.Called(id)
	return args.Get(0).(models.Servicio), args.Error(1)
}
