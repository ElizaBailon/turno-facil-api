package repository_test

import (
	"testing"
	"time"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"

	"github.com/glebarez/sqlite" // 🌟 CAMBIO AQUÍ: Driver puro de Go sin CGO
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqliteTurnosRepo_CrearYBuscar(t *testing.T) {
	// Arrange: Conectamos a una base de datos SQLite limpia en memoria RAM (:memory:)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// Ejecutamos AutoMigrate para crear las tablas necesarias en el entorno de prueba
	err = db.AutoMigrate(&models.Turno{}, &models.Vehiculo{}, &models.Mecanico{}, &models.Servicio{}, &models.Cliente{})
	assert.NoError(t, err)

	// Insertamos registros de soporte mínimos obligatorios para evitar errores de llaves foráneas
	mecanico := models.Mecanico{ID: 1, Nombre: "Juan Mecánico"}
	servicio := models.Servicio{ID: 1, Nombre: "Alineación", DuracionMins: 45}
	db.Create(&mecanico)
	db.Create(&servicio)

	repo := repository.NewTurnosRepository(db)

	nuevoTurno := models.Turno{
		VehiculoID:  1,
		MecanicoID:  1,
		ServicioID:  1,
		FechaHora:   time.Now().Add(24 * time.Hour),
		DuracionEst: 45,
		Estado:      "pendiente",
	}

	// Act: Pasito 1 - Crear el registro usando el repositorio real con GORM
	err = repo.Crear(&nuevoTurno)
	assert.NoError(t, err)
	assert.NotZero(t, nuevoTurno.ID, "GORM debió asignar un ID autoincremental")

	// Act: Pasito 2 - Buscar el registro guardado para verificar que se persista correctamente
	turnoEncontrado, err := repo.BuscarPorID(nuevoTurno.ID)

	// Assert: Validamos que la base de datos refleje exactamente lo que guardamos
	assert.NoError(t, err)
	assert.Equal(t, nuevoTurno.ID, turnoEncontrado.ID)
	assert.Equal(t, "pendiente", turnoEncontrado.Estado)
	assert.Equal(t, 1, turnoEncontrado.MecanicoID)
}
