package repository

import (
	"turno-facil-api/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockMecanicosRepository struct {
	mock.Mock
}

func (m *MockMecanicosRepository) Crear(mecanico *models.Mecanico) error {
	args := m.Called(mecanico)
	return args.Error(0)
}

func (m *MockMecanicosRepository) Listar() ([]models.Mecanico, error) {
	args := m.Called()
	return args.Get(0).([]models.Mecanico), args.Error(1)
}

func (m *MockMecanicosRepository) ObtenerPorID(id int) (models.Mecanico, error) {
	args := m.Called(id)
	return args.Get(0).(models.Mecanico), args.Error(1)
}

// 🌟 AGREGADO: Ahora el mock sí sabe cómo simular la actualización en SQLite
func (m *MockMecanicosRepository) CambiarDisponibilidad(id int, activo bool) error {
	args := m.Called(id, activo)
	return args.Error(0)
}
