package services_test

import (
	"testing"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
	"turno-facil-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test de Regla de Negocio Real: Rechazar datos inválidos (Duración <= 0)
func TestRegistrarServicio_Error_DuracionInvalida(t *testing.T) {
	// Arrange: Preparamos el entorno simulado con el Mock
	mockRepo := new(repository.MockServiciosRepository)
	service := services.NewServiciosService(mockRepo)

	servicioInvalido := models.Servicio{
		Nombre:       "Cambio de Pastillas",
		DuracionMins: 0, // 🌟 Inválido según tu regla: serv.DuracionMins <= 0
	}

	// Act: Invocamos tu función real 'RegistrarServicio'
	_, err := service.RegistrarServicio(servicioInvalido)

	// Assert: Validamos que se ejecute el rechazo estricto de negocio
	assert.Error(t, err)
	assert.Equal(t, "la duración estimada debe ser mayor a 0 minutos", err.Error())

	// Verificamos de forma estricta que el dato NUNCA llegó al repositorio
	mockRepo.AssertNotCalled(t, "Crear", mock.Anything)
}
