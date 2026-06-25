package repository

import (
	"turno-facil-api/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockTurnosRepository es una estructura que simula el comportamiento del repositorio real
type MockTurnosRepository struct {
	mock.Mock
}

func (m *MockTurnosRepository) Crear(t *models.Turno) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTurnosRepository) Listar() ([]models.Turno, error) {
	args := m.Called()
	return args.Get(0).([]models.Turno), args.Error(1)
}

func (m *MockTurnosRepository) ListarPorMecanico(mecanicoID int) ([]models.Turno, error) {
	args := m.Called(mecanicoID)
	return args.Get(0).([]models.Turno), args.Error(1)
}

func (m *MockTurnosRepository) Eliminar(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTurnosRepository) BuscarPorID(id int) (models.Turno, error) {
	args := m.Called(id)
	return args.Get(0).(models.Turno), args.Error(1)
}

func (m *MockTurnosRepository) Actualizar(t *models.Turno) error {
	args := m.Called(t)
	return args.Error(0)
}
