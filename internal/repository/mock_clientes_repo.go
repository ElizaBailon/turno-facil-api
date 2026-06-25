package repository

import (
	"turno-facil-api/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockClientesRepository simula el comportamiento de la base de datos para los tests
type MockClientesRepository struct {
	mock.Mock
}

func (m *MockClientesRepository) Crear(c *models.Cliente) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockClientesRepository) BuscarPorCedula(cedula string) (models.Cliente, error) {
	args := m.Called(cedula)
	// Retornamos el objeto casteado al modelo correspondiente y su error
	return args.Get(0).(models.Cliente), args.Error(1)
}

func (m *MockClientesRepository) Actualizar(c *models.Cliente) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockClientesRepository) Eliminar(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
