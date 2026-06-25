package storage

import (
	"fmt"
	"log"
	"os"
	"turno-facil-api/internal/models"

	"github.com/glebarez/sqlite" // Sigue sirviendo para desarrollo local
	"gorm.io/driver/postgres"    // Requerido para producción/Docker
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB inicializa la base de datos de manera dinámica (Postgres para producción, SQLite para desarrollo)
func InitDB() *gorm.DB {
	var err error

	// Intentamos leer las variables de entorno de PostgreSQL que configuramos en Docker
	dbHost := os.Getenv("DB_HOST")

	if dbHost != "" {
		// 🔥 MODO PRODUCCIÓN: Si existe DB_HOST, nos conectamos a PostgreSQL usando un DSN estructurado
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Guayaquil",
			dbHost, dbUser, dbPass, dbName, dbPort)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Error al conectar a PostgreSQL en producción: ", err)
		}
		log.Println("🚀 Conexión a PostgreSQL (Docker) establecida correctamente.")
	} else {
		// 💻 MODO DESARROLLO LOCAL: Si no hay variables de entorno, cae automáticamente a SQLite
		dbPath := os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "turno-facil.db" // Fallback local estándar
		}

		DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			log.Fatal("Error al conectar a SQLite local: ", err)
		}
		log.Println("💻 Conexión a SQLite (Desarrollo local) establecida correctamente.")
	}

	// AutoMigrate se ejecuta idéntico sin importar el motor gracias a la abstracción de GORM
	err = DB.AutoMigrate(
		&models.Especialidad{},
		&models.Mecanico{},
		&models.Cliente{},
		&models.Vehiculo{},
		&models.Servicio{},
		&models.Turno{},
		&models.Usuario{},
	)
	if err != nil {
		log.Fatal("Error al ejecutar las migraciones: ", err)
	}

	log.Println("✅ Migraciones de GORM completadas de forma exitosa.")
	return DB
}
