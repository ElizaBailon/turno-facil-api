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

// Test 1: Éxito al registrar un nuevo cliente
func TestRegistrarCliente_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	clienteNuevo := &models.Cliente{
		Nombre:   "María Mero",
		Cedula:   "1312345678",
		Vehiculo: []models.Vehiculo{}, // 🌟 Cambiado a Vehiculo
	}

	mockRepo.On("Crear", mock.Anything).Return(nil)

	// Act
	err := service.RegistrarCliente(clienteNuevo)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test 2: Error al registrar cliente con datos vacíos
func TestRegistrarCliente_Error_DatosVacios(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	clienteInvalido := &models.Cliente{
		Nombre: "",
		Cedula: "",
	}

	// Act
	err := service.RegistrarCliente(clienteInvalido)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "el nombre y la cédula son campos obligatorios", err.Error())
}

// Test 3: Error al registrar cliente con cédula con longitud incorrecta
func TestRegistrarCliente_Error_CedulaInvalida(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	clienteInvalido := &models.Cliente{
		Nombre: "Juan Perez",
		Cedula: "1312", // Corta (No tiene 10 dígitos)
	}

	// Act
	err := service.RegistrarCliente(clienteInvalido)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "la cédula debe tener exactamente 10 dígitos", err.Error())
}

// Test 4: Éxito al buscar un cliente por su cédula (Validando que retorne sus vehículos)
func TestObtenerPorCedula_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	clienteFalso := models.Cliente{
		ID:     1,
		Nombre: "Pedro Delgado",
		Cedula: "1309876543",
		Vehiculo: []models.Vehiculo{ // 🌟 Cambiado a Vehiculo
			{ID: 1, Placa: "MBA-1234", Marca: "Chevrolet", ClienteID: 1},
		},
	}

	mockRepo.On("BuscarPorCedula", "1309876543").Return(clienteFalso, nil)

	// Act
	resultado, err := service.ObtenerPorCedula("1309876543")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Pedro Delgado", resultado.Nombre)
	assert.Len(t, resultado.Vehiculo, 1) // 🌟 Cambiado a Vehiculo
	assert.Equal(t, "MBA-1234", resultado.Vehiculo[0].Placa)
	mockRepo.AssertExpectations(t)
}

// Test 5: Error al buscar un cliente que no existe
func TestObtenerPorCedula_Error_NoEncontrado(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	mockRepo.On("BuscarPorCedula", "9999999999").Return(models.Cliente{}, errors.New("record not found"))

	// Act
	_, err := service.ObtenerPorCedula("9999999999")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "record not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// Test 6: Éxito al actualizar un cliente y sus vehículos (Cascada)
func TestActualizarCliente_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	clienteEditado := &models.Cliente{
		ID:     1,
		Nombre: "María Mero Editado",
		Cedula: "1312345678",
	}

	mockRepo.On("Actualizar", clienteEditado).Return(nil)

	// Act
	err := service.Actualizar(clienteEditado)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test 7: Error al actualizar si se mandan datos vacíos
func TestActualizarCliente_Error_CamposVacios(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	clienteInvalido := &models.Cliente{
		ID:     1,
		Nombre: "", // Nombre vacío inválido
		Cedula: "1312345678",
	}

	// Act
	err := service.Actualizar(clienteInvalido)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "el nombre y la cédula no pueden estar vacíos", err.Error())
}

// Test 8: Camino Feliz - Éxito al eliminar un cliente por su ID en string
func TestEliminarCliente_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	mockRepo.On("Eliminar", 1).Return(nil)

	// Act
	err := service.Eliminar("1")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Test 9: Caso de Error - ID en formato inválido o negativo
func TestEliminarCliente_Error_IDInvalido(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	errNonNumeric := service.Eliminar("abc")
	assert.Error(t, errNonNumeric)
	assert.Equal(t, "ID de cliente inválido", errNonNumeric.Error())

	errNegative := service.Eliminar("-5")
	assert.Error(t, errNegative)
	assert.Equal(t, "ID de cliente inválido", errNegative.Error())

	mockRepo.AssertNotCalled(t, "Eliminar", mock.Anything)
}

// Test 10: Caso de Error - El repositorio retorna que no existe el cliente
func TestEliminarCliente_Error_NoEncontrado(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockClientesRepository)
	service := services.NewClientesService(mockRepo)

	mockRepo.On("Eliminar", 99).Return(errors.New("cliente no encontrado"))

	// Act
	err := service.Eliminar("99")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "cliente no encontrado", err.Error())
	mockRepo.AssertExpectations(t)
}
