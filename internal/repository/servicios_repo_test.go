package repository_test

import (
	"testing"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqliteServiciosRepo_CrearYBuscar(t *testing.T) {
	// Arrange: SQLite puro en memoria RAM sin CGO
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// Ejecutar la migración automática de la entidad Servicio
	err = db.AutoMigrate(&models.Servicio{})
	assert.NoError(t, err)

	repo := repository.NewServiciosRepository(db) // Asegúrate de tener este constructor en tu repo de servicios
	nuevoServicio := models.Servicio{
		Nombre:       "Mantenimiento General",
		DuracionMins: 90,
	}

	// Act: Pasito 1 - Crear en base de datos
	err = repo.Crear(&nuevoServicio) // Asumiendo que tu método se llama Crear
	assert.NoError(t, err)
	assert.NotZero(t, nuevoServicio.ID)

	// Act: Pasito 2 - Buscar por ID para reflejar la persistencia
	servicioEncontrado, err := repo.ObtenerPorID(nuevoServicio.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, nuevoServicio.ID, servicioEncontrado.ID)
	assert.Equal(t, "Mantenimiento General", servicioEncontrado.Nombre)
	assert.Equal(t, 90, servicioEncontrado.DuracionMins)
}
