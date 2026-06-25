package main

import (
	"net/http"
	"turno-facil-api/internal/handlers"
	"turno-facil-api/internal/repository"
	"turno-facil-api/internal/services"
	"turno-facil-api/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 1. Inicializar Base de Datos (Postgres en producción / SQLite local)
	db := storage.InitDB()

	// 2. Inicializar Repositorios (Interfaces de acceso a datos)
	clientesRepo := repository.NewClientesRepository(db)
	mecanicosRepo := repository.NewMecanicosRepository(db)
	turnosRepo := repository.NewTurnosRepository(db)
	serviciosRepo := repository.NewServiciosRepository(db)
	usuariosRepo := repository.NewUsuariosRepository(db)

	// 3. Inicializar Servicios (Lógica de Negocio Obligatoria)
	// 🌟 INTEGRANTE B (Tú)
	clientesServ := services.NewClientesService(clientesRepo)

	// 🌟 INTEGRANTE A (Tu compañero)
	mecanicosServ := services.NewMecanicosService(mecanicosRepo)

	// Otros servicios compartidos/grupales
	turnosServ := services.NewTurnosService(turnosRepo, serviciosRepo)
	authServ := services.NewAuthService(usuariosRepo)

	// 🌟 CORRECCIÓN EXIGIDA: Agregamos la capa de negocio para Servicios del Taller
	serviciosServ := services.NewServiciosService(serviciosRepo)

	// 4. Inicializar Controladores (Handlers)
	// 🌟 INTEGRANTE B (Tú): Handler recibe el Servicio
	clientesHand := handlers.NewClientesHandler(clientesServ)

	// 🌟 INTEGRANTE A (Tu compañero): Handler recibe el Servicio
	mecanicosHand := handlers.NewMecanicosHandler(mecanicosServ)

	// 🌟 CORRECCIÓN EXIGIDA: Ahora el Handler recibe el Servicio correspondiente, no el Repositorio
	serviciosHand := handlers.NewServiciosHandler(serviciosServ)

	turnosHand := handlers.NewTurnosHandler(turnosServ)
	authHand := handlers.NewAuthHandler(authServ)

	// Ruta raíz informativa
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"proyecto": "TurnoFácil API", "estado": "operacional", "hito": 3, "arquitectura": "Handler->Service->Repository"}`))
	})

	// 5. Endpoints de Autenticación (Públicos)
	r.Post("/api/v1/auth/register", authHand.Registrar)
	r.Post("/api/v1/auth/login", authHand.Login)

	// 6. Montar las rutas en el Router protegido de Chi
	r.Group(func(protected chi.Router) {
		protected.Use(handlers.AuthMiddleware)

		protected.Mount("/api/v1/mecanicos", mecanicosHand.Routes())
		protected.Mount("/api/v1/servicios", serviciosHand.Routes())
		protected.Mount("/api/v1/clientes", clientesHand.Routes())

		protected.Mount("/api/v1/turnos", turnosHand.Routes())
	})

	println("Servidor de TurnoFácil ejecutándose en http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
