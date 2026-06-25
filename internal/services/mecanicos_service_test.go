package services_test

import (
	"errors"
	"testing"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
	"turno-facil-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test 1: Listar Mecánicos con éxito
func TestListarMecanicos_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockMecanicosRepository)
	service := services.NewMecanicosService(mockRepo)

	// 🌟 CORRECCIÓN: Inicializamos la Especialidad como el objeto que es
	listaFalsa := []models.Mecanico{
		{
			ID:             1,
			Nombre:         "Carlos Zambrano",
			EspecialidadID: 1,
			Especialidad:   models.Especialidad{ID: 1, Nombre: "Motores"},
			Activo:         true,
		},
		{
			ID:             2,
			Nombre:         "Christian Cevallos",
			EspecialidadID: 2,
			Especialidad:   models.Especialidad{ID: 2, Nombre: "Frenos"},
			Activo:         true,
		},
	}

	mockRepo.On("Listar").Return(listaFalsa, nil)

	// Act
	mecanicos, err := service.ListarMecanicos()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, mecanicos, 2)
	assert.Equal(t, "Carlos Zambrano", mecanicos[0].Nombre)
	mockRepo.AssertExpectations(t)
}

// Test 2: Caso de Error - Buscar mecánico que no existe en el taller
func TestObtenerMecanico_Error_NoEncontrado(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockMecanicosRepository)
	service := services.NewMecanicosService(mockRepo)

	mockRepo.On("ObtenerPorID", 99).Return(models.Mecanico{}, errors.New("record not found"))

	// Act
	_, err := service.ObtenerMecanicoPorID(99)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "record not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test 3: Éxito al registrar un nuevo mecánico en el taller
func TestRegistrarMecanico_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockMecanicosRepository)
	service := services.NewMecanicosService(mockRepo)

	// 🌟 CORRECCIÓN: Ajustado al modelo relacional
	mecanicoNuevo := models.Mecanico{
		Nombre:         "Jose Macias",
		EspecialidadID: 3,
		Especialidad:   models.Especialidad{ID: 3, Nombre: "Alineación y Balanceo"},
		Activo:         true,
	}

	mockRepo.On("Crear", mock.Anything).Return(nil)

	// Act
	resultado, err := service.RegistrarMecanico(mecanicoNuevo)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Jose Macias", resultado.Nombre)
	mockRepo.AssertExpectations(t)
}

// Test 4: Error al registrar mecánico con nombre vacío
func TestRegistrarMecanico_Error_NombreVacio(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockMecanicosRepository)
	service := services.NewMecanicosService(mockRepo)

	// 🌟 CORRECCIÓN: Ajustado al modelo relacional
	mecanicoInvalido := models.Mecanico{
		Nombre:         "",
		EspecialidadID: 2,
		Especialidad:   models.Especialidad{ID: 2, Nombre: "Frenos"},
	}

	// Act
	_, err := service.RegistrarMecanico(mecanicoInvalido)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "el nombre del mecánico es obligatorio")
}

// Test 5: Éxito al actualizar la disponibilidad de un mecánico
func TestActualizarEstadoMecanico_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockMecanicosRepository)
	service := services.NewMecanicosService(mockRepo)

	mockRepo.On("CambiarDisponibilidad", int(1), false).Return(nil)

	// Act
	err := service.CambiarDisponibilidad(1, false)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
