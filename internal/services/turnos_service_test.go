package services_test

import (
	"errors"
	"testing"
	"time"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
	"turno-facil-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test 1: Camino Feliz - El turno se registra con éxito porque el mecánico está libre
func TestRegistrarTurno_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockTurnosRepository)
	mockServiciosRepo := new(repository.MockServiciosRepository) // 🌟 NUEVO MOCK

	// 🌟 CORRECCIÓN: Inyectamos ambos mocks al constructor
	service := services.NewTurnosService(mockRepo, mockServiciosRepo)

	fechaTurno := time.Now().Add(24 * time.Hour)

	nuevoTurno := models.Turno{
		VehiculoID: 1,
		MecanicoID: 1,
		ServicioID: 1, // ID del servicio
		FechaHora:  fechaTurno,
		Estado:     "pendiente",
	}

	// 🌟 Simulamos que el servicio mecánico existe y dura 60 minutos
	servicioSimulado := models.Servicio{ID: 1, Nombre: "Cambio de Aceite", DuracionMins: 60}
	mockServiciosRepo.On("ObtenerPorID", 1).Return(servicioSimulado, nil)

	// Simulamos que el mecánico no tiene turnos previos agendados
	mockRepo.On("ListarPorMecanico", 1).Return([]models.Turno{}, nil)
	mockRepo.On("Crear", mock.Anything).Return(nil)

	// Act
	resultado, err := service.RegistrarTurno(nuevoTurno)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, resultado.MecanicoID)
	assert.Equal(t, 60, resultado.DuracionEst) // Valida que heredó la duración
	mockRepo.AssertExpectations(t)
	mockServiciosRepo.AssertExpectations(t)
}

// Test 2: Camino de Error - Falla porque hay un choque exacto de horario
func TestRegistrarTurno_Error_ChoqueHorario(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockTurnosRepository)
	mockServiciosRepo := new(repository.MockServiciosRepository)
	service := services.NewTurnosService(mockRepo, mockServiciosRepo)

	baseTime := time.Now().Add(24 * time.Hour)

	turnosExistentes := []models.Turno{
		{
			VehiculoID:  1,
			MecanicoID:  1,
			FechaHora:   baseTime,
			DuracionEst: 60,
			Estado:      "pendiente",
		},
	}

	turnoConflictivo := models.Turno{
		VehiculoID: 2,
		MecanicoID: 1,
		ServicioID: 1,
		FechaHora:  baseTime,
	}

	// 🌟 El servicio consultado dura 30 minutos
	servicioSimulado := models.Servicio{ID: 1, Nombre: "Revisión Express", DuracionMins: 30}
	mockServiciosRepo.On("ObtenerPorID", 1).Return(servicioSimulado, nil)
	mockRepo.On("ListarPorMecanico", 1).Return(turnosExistentes, nil)

	// Act
	_, err := service.RegistrarTurno(turnoConflictivo)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicto de agenda")
}

// Test 3: Caso Edge - Entrada inválida (El servicio asociado no existe)
func TestRegistrarTurno_Error_ServicioNoExiste(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockTurnosRepository)
	mockServiciosRepo := new(repository.MockServiciosRepository)
	service := services.NewTurnosService(mockRepo, mockServiciosRepo)

	turnoInvalido := models.Turno{
		VehiculoID: 1,
		MecanicoID: 1,
		ServicioID: 999, // Servicio fantasma
		FechaHora:  time.Now().Add(1 * time.Hour),
	}

	// 🌟 Simulamos que el catálogo no encuentra este servicio técnico
	mockServiciosRepo.On("ObtenerPorID", 999).Return(models.Servicio{}, errors.New("gorm: record not found"))

	// Act
	_, err := service.RegistrarTurno(turnoInvalido)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "el servicio seleccionado no existe en el catálogo", err.Error())
}

// Test 4: Caso de Error - Fallo general al intentar guardar en la base de datos
func TestRegistrarTurno_Error_FalloBaseDatos(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockTurnosRepository)
	mockServiciosRepo := new(repository.MockServiciosRepository)
	service := services.NewTurnosService(mockRepo, mockServiciosRepo)

	turnoNormal := models.Turno{
		VehiculoID: 1,
		MecanicoID: 1,
		ServicioID: 1,
		FechaHora:  time.Now().Add(24 * time.Hour),
	}

	servicioSimulado := models.Servicio{ID: 1, Nombre: "Alineación", DuracionMins: 45}
	mockServiciosRepo.On("ObtenerPorID", 1).Return(servicioSimulado, nil)
	mockRepo.On("ListarPorMecanico", 1).Return([]models.Turno{}, nil)
	mockRepo.On("Crear", mock.Anything).Return(errors.New("error interno de base de datos"))

	// Act
	_, err := service.RegistrarTurno(turnoNormal)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error interno de base de datos")
}

func TestActualizarTurno_Exito(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockTurnosRepository)
	mockServiciosRepo := new(repository.MockServiciosRepository)
	service := services.NewTurnosService(mockRepo, mockServiciosRepo)

	idTurno := 1
	turnoExistente := models.Turno{ID: idTurno, MecanicoID: 1, DuracionEst: 60}
	turnoActualizado := models.Turno{MecanicoID: 1, ServicioID: 2, FechaHora: time.Now()}

	mockRepo.On("BuscarPorID", idTurno).Return(turnoExistente, nil)
	mockServiciosRepo.On("ObtenerPorID", 2).Return(models.Servicio{ID: 2, DuracionMins: 30}, nil)
	mockRepo.On("ListarPorMecanico", 1).Return([]models.Turno{}, nil)
	mockRepo.On("Actualizar", mock.Anything).Return(nil)
	mockRepo.On("BuscarPorID", idTurno).Return(turnoExistente, nil) // Segunda llamada tras el update

	// Act
	_, err := service.Actualizar(idTurno, turnoActualizado)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
